package operations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CheckIfExists(collection *mongo.Collection, filter bson.M) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Find(ctx, filter)
	if err != nil {
		return false
	}
	return true
}
