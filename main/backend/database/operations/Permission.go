package operations

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/utils"
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

func GetPermissionByName(db *mongo.Database, permissionName string) (models.Permission, error) {
	var permission models.Permission
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("permission").FindOne(ctx, bson.M{"name": permissionName}).Decode(&permission)

	return permission, err
}

func GetAllPermissionsFromRole(db *mongo.Database, roleId primitive.ObjectID) ([]models.Permission, error) {
	var permissions []models.Permission

	role, err := GetRoleById(db, roleId)
	if err != nil {
		return permissions, err
	}

	for _, permissionId := range role.Permissions {
		permission, err := GetPermissionById(db, permissionId)
		if err != nil {
			return permissions, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func CreatePermission(name string, description string) (models.Permission, error) {
	var permission models.Permission

	permission.ID = primitive.NewObjectID()
	permission.Name = name
	permission.Description = description
	permission.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	permission.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	// check if permission already exists
	existingPermission, err := GetPermissionByName(database.MongoDB, permission.Name)
	if err == nil {
		return existingPermission, nil
	}

	err = SavePermission(database.MongoDB, permission)
	if err != nil {
		return permission, err
	}

	return permission, nil
}
