package middleware

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/crypt"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type JwtCheck struct {
	ErrorCode int
	Message   string
	Abort     bool
}

// JWTAuth is a middleware that checks if the user is logged in
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		check := CheckJWT(c)
		if check.Abort {
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(check.ErrorCode, gin.H{"msg": check.Message, "error": true})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// CheckJWT checks if the JWT token is valid and if the user is logged in
func CheckJWT(c *gin.Context) JwtCheck {
	unauthorizedError := JwtCheck{
		ErrorCode: http.StatusUnauthorized,
		Message:   "Unauthorized",
		Abort:     true,
	}

	// check if the token exists
	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		log.Println("The auth_token is missing in the cookie")
		return unauthorizedError
	}

	if tokenString == "" {
		log.Println("The auth_token is empty")
		return unauthorizedError
	}

	// parse the token to check the claims and the signature
	jwtToken, err := crypt.ParseJwt(tokenString)
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	// get the userId from the jwtToken
	userId, err := primitive.ObjectIDFromHex(jwtToken["userId"].(string))
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	// check if user exists
	user, err := operations.GetOneById[models.User](database.Tables.User, userId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
		}
		return unauthorizedError
	}

	// check if the token is still valid
	exp, ok := jwtToken["exp"].(float64)
	if !ok {
		return unauthorizedError
	}

	if time.Now().Unix() > int64(exp) {
		return unauthorizedError
	}

	// update the context
	c.Set("user", user)
	c.Set("userId", jwtToken["userId"])
	c.Set("loggedIn", true)

	return JwtCheck{
		ErrorCode: 0,
		Message:   "Authorized",
		Abort:     false,
	}
}
