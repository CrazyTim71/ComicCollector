package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	JoiValidator "ComicCollector/main/backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

func RegisterHandler(rg *gin.RouterGroup) {
	rg.POST("/", func(c *gin.Context) {
		var requestBody struct {
			Username         string `json:"username" binding:"required"`
			Password         string `json:"password" binding:"required"`
			PasswordRepeated string `json:"passwordRepeated" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": "true"})
		}

		// check if password and repeated password match
		if requestBody.Password != requestBody.PasswordRepeated {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Passwords do not match", "error": "true"})
		}

		// check if username and password are allowed
		if err := JoiValidator.UsernameSchema.Validate(requestBody.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": "true"})
		}

		if err := JoiValidator.PasswordSchema.Validate(requestBody.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid password. Please follow the password rules.", "error": "true"})
		}

		// check if the user already exists in the database by querying with the username
		var existingUser models.User

		// prevents multiple variations of the same username with uppercase and lowercase letters
		username := strings.ToLower(requestBody.Username)

		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"username": username}).Decode(existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "This username already exists", "error": "true"})
		}
	})
}
