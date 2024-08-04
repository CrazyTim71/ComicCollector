package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SaveUser(db *mongo.Database, newUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("user").InsertOne(ctx, newUser, options.InsertOne())
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(db *mongo.Database, username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := db.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)

	return existingUser, err
}

func GetUserById(db *mongo.Database, id primitive.ObjectID) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := db.Collection("user").FindOne(ctx, bson.M{"_id": id}).Decode(&existingUser)

	return existingUser, err
}
