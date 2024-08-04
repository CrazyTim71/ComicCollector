package middleware

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/crypt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			log.Println(err)
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

func VerifyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the user is logged in
		loggedIn, exists := c.Get("loggedIn")
		if !exists || !loggedIn.(bool) {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// check if the user is an admin
		userId, exists := c.Get("userId")
		if !exists {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		user, err := operations.GetUserById(database.MongoDB, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// get all roles by the user id
		userRoles, err := operations.GetUserRolesByUserId(database.MongoDB, user.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if the user has the admin role
		isAdmin := false
		for _, userRole := range userRoles {
			// get single role
			role, err := operations.GetRoleById(database.MongoDB, userRole.RoleId)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
				c.Abort()
				return
			}

			if role.Name == "Administrator" {
				isAdmin = true
				break
			}
		}
		if isAdmin {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "Not enough permissions to view this site", "error": true})
			c.Abort()
			return
		}

	}
}
