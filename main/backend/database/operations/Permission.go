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

func SavePermission(db *mongo.Database, newPermission models.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("permission").InsertOne(ctx, newPermission, options.InsertOne())

	return err
}

func GetPermissionById(db *mongo.Database, permissionId primitive.ObjectID) (models.Permission, error) {
	var permission models.Permission
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("permission").FindOne(ctx, bson.M{"_id": permissionId}).Decode(&permission)

	return permission, err
}
