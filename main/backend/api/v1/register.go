package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/helpers"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/JoiHelper"
	"ComicCollector/main/backend/utils/crypt"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
	"time"
)

// RegisterHandler api/v1/register
func RegisterHandler(rg *gin.RouterGroup) {
	rg.GET("check", func(c *gin.Context) {
		if env.GetSignupEnabled() {
			c.JSON(http.StatusOK, gin.H{"signupEnabled": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"signupEnabled": false})
		}
	})

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
		if err := JoiHelper.UsernameSchema.Validate(requestBody.Username); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": true})
			return
		}

		if err := JoiHelper.PasswordSchema.Validate(requestBody.Password); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid password. Please follow the password rules.", "error": true})
			return
		}

		// prevents multiple variations of the same username with uppercase and lowercase letters
		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password

		// check if the user already exists in the database by querying with the username
		_, err := operations.GetOneByFilter[models.User](database.Tables.User, bson.M{"username": username})
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

		// create the restricted user permissions
		var restrictedUserPermissionIds []primitive.ObjectID
		for _, permission := range groups.RestrictedUser.Permissions {
			perm, err := helpers.CreatePermission(permission.Name(), permission.Description())
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			restrictedUserPermissionIds = append(restrictedUserPermissionIds, perm.ID)
		}

		// create the roles in case they don't exist
		restrictedUserRole, err := helpers.CreateRole(
			database.MongoDB,
			groups.RestrictedUser.Name,
			groups.RestrictedUser.Description,
			restrictedUserPermissionIds,
		)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			return
		}

		// add the role to the user
		newUser.Roles = append(newUser.Roles, restrictedUserRole.ID)

		_, err = operations.InsertOne(database.Tables.User, newUser)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			return
		}

		c.Redirect(http.StatusSeeOther, "/login")
	})
}
