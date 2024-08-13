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

func GetAllUsers(db *mongo.Database) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("user").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = cursor.All(ctx, &users)

	// remove the password hash because we don't want to expose that xD
	for i := range users {
		users[i].Password = ""
	}

	return users, err
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

	// remove the password hash because we don't want to expose that xD
	existingUser.Password = ""

	return existingUser, err
}

func CreateUser(db *mongo.Database, newUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("user").InsertOne(ctx, newUser, options.InsertOne())

	return err
}

func UpdateUserById(db *mongo.Database, id primitive.ObjectID, updatedUser models.User) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("user").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedUser})

	return result, err
}

func DeleteUserById(db *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("user").DeleteOne(ctx, bson.M{"_id": id})

	return err
}
