package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/utils/Joi"
	"ComicCollector/main/backend/utils/crypt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
)

func RegisterHandler(rg *gin.RouterGroup) {
	rg.POST("", func(c *gin.Context) {
		var requestBody struct {
			Username         string `json:"username" binding:"required"`
			Password         string `json:"password" binding:"required"`
			PasswordRepeated string `json:"passwordRepeated" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
			log.Println(err)
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

		// check if the user already exists in the database by querying with the username
		var existingUser models.User

		// prevents multiple variations of the same username with uppercase and lowercase letters
		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password

		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"username": username}).Decode(&existingUser)
		if err == nil { // err == nil in case the user already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "This username already exists", "error": true})
			log.Println(err)
			return
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
		newUser.ID = primitive.NewObjectID()
		newUser.Username = username
		newUser.Password = hashedPW
		newUser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		newUser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = database.MongoDB.Collection("user").InsertOne(c, newUser, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
			log.Println(err)
			return
		}

		// TODO: redirect to the login
		// c.Redirect(http.StatusSeeOther, "/login")
		c.JSON(http.StatusOK, gin.H{"msg": "User was created successfully"})
	})
}
