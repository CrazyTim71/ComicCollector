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

func SaveRolePermission(db *mongo.Database, newRolePermission models.RolePermission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("role_permission").InsertOne(ctx, newRolePermission, options.InsertOne())

	return err
}

func GetRolePermissionById(db *mongo.Database, rolePermissionId primitive.ObjectID) (models.RolePermission, error) {
	var rolePermission models.RolePermission
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("role_permission").FindOne(ctx, bson.M{"_id": rolePermissionId}).Decode(&rolePermission)

	return rolePermission, err
}

func GetAllRolePermissionsByRoleId(db *mongo.Database, roleId primitive.ObjectID) ([]models.RolePermission, error) {
	var rolePermissions []models.RolePermission
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("role_permission").Find(ctx, bson.M{"role_id": roleId})
	if err != nil {
		return rolePermissions, err
	}

	err = cursor.All(ctx, &rolePermissions)

	return rolePermissions, err
}
