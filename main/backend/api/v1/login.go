package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/Joi"
	"ComicCollector/main/backend/utils/crypt"
	"github.com/gin-gonic/gin"
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

		// check username input
		if err := Joi.UsernameSchema.Validate(requestBody.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": true})
			log.Println(err)
			return
		}

		// check password input
		if err := Joi.PasswordSchema.Validate(requestBody.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid password. Please remove all invalid characters and try again.", "error": true})
			log.Println(err)
			return
		}

		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password

		// check if the user exists
		existingUser, err := operations.GetUserByUsername(database.MongoDB, username)
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
