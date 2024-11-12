package middleware

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/crypt/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type JwtCheck struct {
	ErrorCode int
	Message   string
	Abort     bool
	UserId    primitive.ObjectID
}

var unauthorizedError = JwtCheck{
	ErrorCode: http.StatusUnauthorized,
	Message:   "Unauthorized",
	Abort:     true,
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
	refreshCheck := checkJWTRefresh(c)

	// Refresh expired
	if refreshCheck.Abort {
		auth.DeleteAuthCookies(c)
		return refreshCheck
	}

	authCheck := checkJWTAuth(c)

	// Auth expired, but refresh token is valid
	if authCheck.Abort {
		// Request new access token, because refresh token is valid
		err := auth.GenerateTokens(c, refreshCheck.UserId)
		if err != nil {
			log.Println(err)
			return unauthorizedError
		}
		err = setContextUser(c, refreshCheck.UserId)
		if err != nil {
			log.Println(err)
			return unauthorizedError
		}
		return JwtCheck{
			ErrorCode: 0,
			Message:   "Authorized",
			Abort:     false,
		}
	}

	return authCheck

}

func checkJWTAuth(c *gin.Context) JwtCheck {
	tokenString, err := c.Cookie(auth.AuthCookieName)
	if err != nil {
		log.Println("The auth_token is missing in the cookie")
		return unauthorizedError
	}

	if tokenString == "" {
		log.Println("The auth_token is empty")
		return unauthorizedError
	}

	jwtToken, err := auth.ParseJwt(tokenString)
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	id, err := primitive.ObjectIDFromHex(jwtToken["userId"].(string))
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	err = setContextUser(c, id)
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	return JwtCheck{
		ErrorCode: 0,
		Message:   "Authorized",
		Abort:     false,
		UserId:    id,
	}
}

func setContextUser(c *gin.Context, userId primitive.ObjectID) error {
	dUser, err := operations.GetOneById[models.User](database.Tables.User, userId)
	if err != nil {
		log.Println(err)
		return err
	}

	c.Set("user", dUser)
	c.Set("userId", userId)
	c.Set("loggedIn", true)

	return nil
}

func checkJWTRefresh(c *gin.Context) JwtCheck {
	tokenString, err := c.Cookie(auth.AuthCookieRefreshName)

	if err != nil {
		log.Println("The auth_refresh_token is missing in the cookie")
		return unauthorizedError
	}

	if tokenString == "" {
		log.Println("The auth_refresh_token is empty")
		return unauthorizedError
	}

	jwtToken, err := auth.ParseJwt(tokenString)
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	tokenClaim := jwtToken["token"].(string)

	refreshToken, err := operations.GetOneByFilter[models.AuthToken](database.Tables.AuthToken, bson.M{"token": tokenClaim})
	if err != nil {
		log.Println(err)
		return unauthorizedError
	}

	if refreshToken.Type != models.AuthRefreshToken {
		log.Println("The token is not a refresh token")
		return unauthorizedError
	}

	if auth.IsExpired(refreshToken.ExpiresAt) {
		log.Println("The refresh token has expired")
		_, _ = operations.DeleteOne(database.Tables.AuthToken, bson.M{"_id": refreshToken.ID})
		return unauthorizedError
	}

	return JwtCheck{
		ErrorCode: 0,
		Message:   "Authorized",
		Abort:     false,
		UserId:    refreshToken.UserId,
	}
}
