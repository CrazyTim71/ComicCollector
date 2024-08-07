package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils/Joi"
	"ComicCollector/main/backend/utils/crypt"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
	"time"
)

// UserHandler api/v1/user
func UserHandler(rg *gin.RouterGroup) {
	// rg.POST() to create a new user is done in register.go

	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.VerifyHasAllPermission(
			permissions.UserViewAll,
		),
		func(c *gin.Context) {
			// returns all users
			users, err := operations.GetAllUsers(database.MongoDB)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			}
			c.JSON(http.StatusOK, users)
		})

	rg.GET("/:id",
		middleware.CheckJwtToken(),
		func(c *gin.Context) {
			id := c.Param("id")

			userId, exits := c.Get("userId")
			if !exits {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
				return
			}

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the user has the permission to view himself
			if userId == id {
				middleware.VerifyHasAllPermission(
					permissions.UserViewSelf,
				)(c)

				// otherwise check if the user has admin rights
			} else {
				middleware.VerifyHasAllPermission(
					permissions.UserViewAll,
				)(c)
			}

			// abort if the permission check failed
			if c.IsAborted() {
				return
			}

			user, err := operations.GetUserById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			}

			c.JSON(http.StatusOK, user)
		})

	rg.PATCH("/:id",
		middleware.CheckJwtToken(),
		func(c *gin.Context) {
			{
				id := c.Param("id")

				userId, exits := c.Get("userId")
				if !exits {
					c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
					return
				}

				objID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
					return
				}

				// check if the user has the permission to modify himself
				if userId == id {
					middleware.VerifyHasAllPermission(
						permissions.UserModifySelf,
					)(c)

					// otherwise check if the user has admin rights
				} else {
					middleware.VerifyHasAllPermission(
						permissions.UserModifyAll,
					)(c)
				}

				// abort if the permission check failed
				if c.IsAborted() {
					return
				}

				var requestBody struct {
					Username         string `json:"username" binding:"required"`
					Password         string `json:"password" binding:"required"`
					PasswordRepeated string `json:"passwordRepeated" binding:"required"`
				}

				if err := c.ShouldBindJSON(&requestBody); err != nil {
					log.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
					return
				}

				// check if password and repeated password match
				if requestBody.Password != requestBody.PasswordRepeated {
					c.JSON(http.StatusBadRequest, gin.H{"msg": "Passwords do not match", "error": true})
					return
				}

				// check if username and password are allowed
				if err := Joi.UsernameSchema.Validate(requestBody.Username); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": true})
					log.Println(err)
					return
				}

				if err := Joi.PasswordSchema.Validate(requestBody.Password); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid password. Please follow the password rules.", "error": true})
					log.Println(err)
					return
				}

				// prevents multiple variations of the same username with uppercase and lowercase letters
				username := strings.ToLower(requestBody.Username)
				password := requestBody.Password

				// check if the username was changed
				usernameChanged := false
				existingUser, err := operations.GetUserById(database.MongoDB, objID)
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
					return
				}
				if existingUser.Username != username {
					usernameChanged = true
				}

				// check if the user already exists in the database by querying with the username
				_, err = operations.GetUserByUsername(database.MongoDB, username)
				if err == nil { // err == nil in case the user already exists
					// only show the error when the username was changed
					// otherwise its logical that the user already exists
					if usernameChanged {
						c.JSON(http.StatusConflict, gin.H{"msg": "This username already exists", "error": true})
						log.Println(err)
						return
					}
				} else if !errors.Is(err, mongo.ErrNoDocuments) {
					// handle all other database errors, but ignore the NoDocuments error
					// that's because this error is expected when the user doesn't exist
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal database error", "error": true})
					log.Println(err)
					return
				}

				// hash the password
				hashedPW, err := crypt.HashPassword(password)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": true})
					log.Println(err)
					return
				}

				var newUser models.User
				newUser.ID = existingUser.ID
				newUser.Username = username
				newUser.Password = hashedPW
				newUser.CreatedAt = existingUser.CreatedAt
				newUser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

				result, err := operations.UpdateUserById(database.MongoDB, objID, newUser)
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
					return
				}

				if result.ModifiedCount == 0 {
					c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing modified"})
					return
				} else {
					c.JSON(http.StatusOK, gin.H{"msg": "Updated user successfully"})
				}
			}
		})

	rg.DELETE("/:id",
		middleware.CheckJwtToken(),
		func(c *gin.Context) {
			{
				id := c.Param("id")

				userId, exits := c.Get("userId")
				if !exits {
					c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
					return
				}

				objID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
					return
				}

				// check if the user has the permission to delete himself
				if userId == id {
					middleware.VerifyHasAllPermission(
						permissions.UserDeleteSelf,
					)(c)

					// otherwise check if the user has admin rights
				} else {
					middleware.VerifyHasAllPermission(
						permissions.UserDeleteAll,
					)(c)
				}

				// abort if the permission check failed
				if c.IsAborted() {
					return
				}

				err = operations.DeleteUserById(database.MongoDB, objID)
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				}

				c.JSON(http.StatusOK, gin.H{"msg": "Deleted user successfully"})
			}
		})
}
