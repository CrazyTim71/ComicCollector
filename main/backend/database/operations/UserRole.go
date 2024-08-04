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

func SaveUserRole(db *mongo.Database, newUserRole models.UserRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("user_role").InsertOne(ctx, newUserRole, options.InsertOne())
	if err != nil {
		return err
	}

	return nil
}

func GetUserRoleById(db *mongo.Database, userRoleId string) (models.UserRole, error) {
	var userRole models.UserRole
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("user_role").FindOne(ctx, bson.M{"_id": userRoleId}).Decode(&userRole)
	if err != nil {
		return userRole, err
	}

	return userRole, nil
}

func GetUserRolesByUserId(db *mongo.Database, userId primitive.ObjectID) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("user_role").Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return userRoles, err
	}

	err = cursor.All(ctx, &userRoles)
	if err != nil {
		return userRoles, err
	}

	return userRoles, nil
}
