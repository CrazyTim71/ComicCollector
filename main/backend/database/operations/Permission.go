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

func GetAllPermissions(db *mongo.Database) ([]models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("permission").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var permissions []models.Permission
	err = cursor.All(ctx, &permissions)

	return permissions, err
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

func CheckIfAllPermissionsExist(db *mongo.Database, permissionIds []primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("permission").Find(ctx, bson.M{"_id": bson.M{"$in": permissionIds}})
	if err != nil {
		return false
	}

	var foundRoles []models.Role
	err = cursor.All(ctx, &foundRoles)

	return len(foundRoles) == len(permissionIds)
}

func InsertPermission(db *mongo.Database, newPermission models.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("permission").InsertOne(ctx, newPermission, options.InsertOne())

	return err
}

func UpdatePermission(db *mongo.Database, permissionId primitive.ObjectID, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("permission").UpdateOne(ctx, bson.M{"_id": permissionId}, bson.M{"$set": data})

	return result, err
}

func DeletePermission(db *mongo.Database, permissionId primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("permission").DeleteOne(ctx, bson.M{"_id": permissionId})

	return result, err
}
