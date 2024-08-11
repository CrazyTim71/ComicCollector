package middleware

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/crypt"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func CheckJwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the token exists
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			//log.Println("The auth_token is missing in the cookie")
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		if tokenString == "" {
			log.Println("The auth_token is empty")
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// parse the token
		jwtToken, err := crypt.ParseJwt(tokenString)
		if err != nil {
			log.Println(err)
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// get the userId from the jwtToken
		userId, err := primitive.ObjectIDFromHex(jwtToken["userId"].(string))
		if err != nil {
			log.Println(err)
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if user exists
		_, err = operations.GetUserById(database.MongoDB, userId)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				log.Println(err)
			}
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if the token is still valid
		if jwtToken["exp"] == nil {
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		if jwtToken["exp"].(float64) < jwtToken["iat"].(float64) {
			// TODO: redirect to login ?
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// update the context
		c.Set("userId", jwtToken["userId"])
		c.Set("loggedIn", true)
		c.Next()
	}
}
