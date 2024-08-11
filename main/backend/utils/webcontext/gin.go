package webcontext

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserId(c *gin.Context) (primitive.ObjectID, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return primitive.NilObjectID, errors.New("userId doesn't exist")
	}

	if userId == "" || userId == nil {
		return primitive.NilObjectID, errors.New("userId is empty")
	}

	id, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		return primitive.NilObjectID, errors.New("userId doesn't exist")
	}

	return id, nil
}
