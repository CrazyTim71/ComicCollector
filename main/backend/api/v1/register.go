package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/Joi"
	"ComicCollector/main/backend/utils/crypt"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
	"time"
)

func RegisterHandler(rg *gin.RouterGroup) {
	rg.POST("", func(c *gin.Context) {
		if !env.GetSignupEnabled() {
			c.JSON(http.StatusForbidden, gin.H{"msg": "Signup is disabled", "error": true})
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
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": true})
			return
		}

		if err := Joi.PasswordSchema.Validate(requestBody.Password); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid password. Please follow the password rules.", "error": true})
			return
		}

		// prevents multiple variations of the same username with uppercase and lowercase letters
		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password

		// check if the user already exists in the database by querying with the username
		_, err := operations.GetUserByUsername(database.MongoDB, username)
		if err == nil { // err == nil in case the user already exists
			log.Println(err)
			c.JSON(http.StatusConflict, gin.H{"msg": "This username already exists", "error": true})
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// handle all other database errors, but ignore the NoDocuments error
			// that's because this error is expected when the user doesn't exist
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal database error", "error": true})
			return
		}

		// hash the password
		hashedPW, err := crypt.HashPassword(password)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": true})
			return
		}

		var newUser models.User
		newUser.ID = primitive.NewObjectID()
		newUser.Username = username
		newUser.Password = hashedPW
		newUser.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
		newUser.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

		err = operations.CreateUser(database.MongoDB, newUser)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			return
		}

		// add the user to the RestrictedUser role for approval
		for _, permission := range groups.RestrictedUser.Permissions {
			_, err := operations.CreatePermission(permission.Name, permission.Description)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
		}

		restrictedUserRole, err := operations.CreateRole(groups.RestrictedUser.Name, groups.RestrictedUser.Description)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			return
		}

		for _, permission := range groups.RestrictedUser.Permissions {
			perm, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			_, err = operations.CreateRolePermission(restrictedUserRole, perm)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
		}

		// assign the user role
		_, err = operations.CreateUserRole(newUser, restrictedUserRole)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			return
		}

		c.Redirect(http.StatusSeeOther, "/login")
		//c.JSON(http.StatusOK, gin.H{"msg": "User was created successfully"})
	})
}
