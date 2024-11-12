package auth

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const (
	AuthCookieName             = "auth_token"
	AuthCookieExpiresIn        = time.Minute * 8
	AuthCookieRefreshName      = "auth_refresh"
	AuthCookieRefreshExpiresIn = time.Hour * 24 * 7
)

var rsaKey *rsa.PrivateKey = nil

func KeySetup(key *rsa.PrivateKey) {
	rsaKey = key
}

func GenerateTokens(c *gin.Context, userId primitive.ObjectID) error {
	accessToken, err := GenerateAccessToken(userId)
	if err != nil {
		log.Println(err)
		return err
	}

	refreshToken, err := GenerateRefreshToken(userId)
	if err != nil {
		log.Println(err)
		return err
	}

	SetAuthCookie(c, accessToken, refreshToken)
	SetRefreshCookie(c, refreshToken)

	return nil
}

func GenerateAccessToken(userId primitive.ObjectID) (string, error) {
	claims := getDefaultClaims(time.Minute * 8)
	claims["userId"] = userId.Hex()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := accessToken.SignedString(rsaKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userId primitive.ObjectID) (string, error) {
	claims := getDefaultClaims(AuthCookieRefreshExpiresIn)
	claims["userId"] = userId.Hex()

	// get all the tokens from a user from the database
	tokens, err := operations.GetManyByFilter[models.AuthToken](database.Tables.AuthToken, bson.M{"userId": userId})
	if err != nil {
		log.Println(err)
		return "", err
	}

	// delete all the refresh tokens from the database to end the sessions
	var refreshTokens []primitive.ObjectID
	for _, t := range tokens {
		if t.Type == models.AuthRefreshToken {
			refreshTokens = append(refreshTokens, t.ID)
		}
	}

	_, _ = operations.DeleteMany(database.Tables.AuthToken, bson.M{"_id": bson.M{"$in": refreshTokens}})

	// generate a unique token
	token := utils.GenerateRandomPassword(16, true, true)
	for {
		tokens, err = operations.GetManyByFilter[models.AuthToken](database.Tables.AuthToken, bson.M{
			"token": token,
		})

		if len(tokens) == 0 {
			break
		}

		token = utils.GenerateRandomPassword(16, true, true)
	}

	claims["token"] = token

	refreshToken := models.AuthToken{
		ID:        primitive.NewObjectID(),
		Token:     token,
		Type:      models.AuthRefreshToken,
		UserId:    userId,
		Expires:   true,
		ExpiresAt: claims["exp"].(int64),
		IssuedAt:  claims["iat"].(int64),
	}

	_, err = operations.InsertOne(database.Tables.AuthToken, refreshToken)
	if err != nil {
		log.Println(err)
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := jwtToken.SignedString(rsaKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

func DeleteAuthCookies(c *gin.Context) {
	c.SetCookie(AuthCookieName, "", 0, "/", "", false, false)
	c.SetCookie(AuthCookieRefreshName, "", 0, "/", "", false, false)
}

func SetAuthCookie(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(AuthCookieName, accessToken, int(AuthCookieExpiresIn.Seconds()), "/", "", true, true)
}

func SetRefreshCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(AuthCookieRefreshName, refreshToken, int(AuthCookieRefreshExpiresIn.Seconds()), "/", "", true, true)
}

func ParseJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("got unexpected jwt signing method")
		}
		return rsaKey.Public(), nil
	}, jwt.WithIssuer("ComicCollector"), jwt.WithExpirationRequired())

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	// Check expiration
	if isExpired(claims) {
		return nil, fmt.Errorf("the token is expired")
	}

	// Check if the token is valid
	if token.Valid {
		return claims, nil
	} else {
		log.Printf("Failed to parse JWT")
		return nil, fmt.Errorf("the token is invalid")
	}
}

// getExpiryDate returns the unix time of the current time plus the duration
func getExpiryDate(duration time.Duration) int64 {
	return time.Now().Add(duration).Unix()
}

// getCurrentTime returns the current unix time
func getCurrentTime() int64 {
	return time.Now().Unix()
}

// IsExpired checks if a unix time is expired
func IsExpired(unixTime int64) bool {
	return time.Now().After(time.Unix(unixTime, 0))
}

// isExpired checks if a JWT token is expired
func isExpired(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(int64)
	if !ok {
		return false
	}
	return IsExpired(exp)
}

func getDefaultClaims(timeDurationValid time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"exp":  getExpiryDate(timeDurationValid),
		"iat":  getCurrentTime(),
		"iss":  "ComicCollector",
		"i am": "your father :O",
	}
}
