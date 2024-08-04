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

func SaveRole(db *mongo.Database, newRole models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("role").InsertOne(ctx, newRole, options.InsertOne())
	if err != nil {
		return err
	}

	return nil
}

func GetRoleById(db *mongo.Database, roleId primitive.ObjectID) (models.Role, error) {
	var role models.Role
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("role").FindOne(ctx, bson.M{"_id": roleId}).Decode(&role)
	if err != nil {
		return role, err
	}

	return role, nil
}
