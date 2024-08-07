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

	return err
}

func GetUserRoleById(db *mongo.Database, userRoleId primitive.ObjectID) (models.UserRole, error) {
	var userRole models.UserRole
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("user_role").FindOne(ctx, bson.M{"_id": userRoleId}).Decode(&userRole)

	return userRole, err
}

func GetUserRoleByName(db *mongo.Database, userRoleName string) (models.UserRole, error) {
	var userRole models.UserRole
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("user_role").FindOne(ctx, bson.M{"name": userRoleName}).Decode(&userRole)

	return userRole, err
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

	return userRoles, err
}
