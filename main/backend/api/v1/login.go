package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/utils/Joi"
	"ComicCollector/main/backend/utils/crypt"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoginHandler(rg *gin.RouterGroup) {
	rg.POST("", func(c *gin.Context) {

		// TODO: make middleware to check the jwt token on every endpoint
		var requestBody struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
			log.Println(err)
			return
		}

		// check if username is allowed
		if err := Joi.UsernameSchema.Validate(requestBody.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": true})
			log.Println(err)
			return
		}

		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password
		var existingUser models.User

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// check if the user exists
		err := database.MongoDB.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials", "error": true})
			log.Println(err)
			return
		} else {
			if crypt.CheckPasswordHash(password, existingUser.Password) {
				// TODO: generate jwt
				// TODO: set cookie
				// TODO: redirect to the dashboard
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials", "error": true})
				log.Println(err)
				return
			}
		}
	})

}
