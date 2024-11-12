package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserId(c *gin.Context) (primitive.ObjectID, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return primitive.NilObjectID, errors.New("userId not found in context")
	}

	return userId.(primitive.ObjectID), nil
}
