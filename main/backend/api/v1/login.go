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
		var requestBody struct {
			Username string `form:"username" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.ShouldBind(&requestBody); err != nil {
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

				// generate a jwt token that will authenticate the user
				jwtToken, err := crypt.GenerateJwtToken(existingUser.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to generate a token", "error": true})
					log.Println(err)
					return
				}

				// save the jwt token as cookie (valid for 24 hours)
				duration := 24 * time.Hour
				cookie := http.Cookie{
					Name:     "auth_token",
					Value:    jwtToken,
					Path:     "/",
					Domain:   "", // Set your domain if needed
					Expires:  time.Now().Add(duration),
					SameSite: http.SameSiteLaxMode,
					Secure:   true, // Ensure you use HTTPS
					HttpOnly: true,
				}
				http.SetCookie(c.Writer, &cookie)

				c.Redirect(http.StatusSeeOther, "/dashboard")
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials", "error": true})
				log.Println(err)
				return
			}
		}
	})
}
