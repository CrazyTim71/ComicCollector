package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/JoiHelper"
	"ComicCollector/main/backend/utils/crypt"
	"ComicCollector/main/backend/utils/crypt/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strings"
)

// LoginHandler api/v1/login
func LoginHandler(rg *gin.RouterGroup) {
	rg.POST("", func(c *gin.Context) {
		var requestBody struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBind(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
			log.Println(err)
			return
		}

		// check username input
		if err := JoiHelper.UsernameSchema.Validate(requestBody.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid username. Please remove all invalid characters and try again.", "error": true})
			log.Println(err)
			return
		}

		// check password input
		if err := JoiHelper.PasswordSchema.Validate(requestBody.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid password. Please remove all invalid characters and try again.", "error": true})
			log.Println(err)
			return
		}

		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password

		// check if the user exists
		existingUser, err := operations.GetOneByFilter[models.User](database.Tables.User, bson.M{"username": username})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials", "error": true})
			log.Println(err)
			return
		} else {
			if crypt.CheckPasswordHash(password, existingUser.Password) {

				// generate a jwt token that will authenticate the user
				err = auth.GenerateTokens(c, existingUser.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to generate a token", "error": true})
					log.Println(err)
					return
				}

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
