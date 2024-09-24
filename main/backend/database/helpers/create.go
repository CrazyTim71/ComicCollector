package helpers

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CreatePermission(name string, description string) (models.Permission, error) {
	var permission models.Permission

	permission.ID = primitive.NewObjectID()
	permission.Name = name
	permission.Description = description
	permission.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	// check if permission already exists
	existingPermission, err := operations.GetOneByFilter[models.Permission](database.Tables.Permission, bson.M{"name": permission.Name})
	if err == nil {
		return existingPermission, nil
	}

	_, err = operations.InsertOne(database.Tables.Permission, permission)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func CreateRole(db *mongo.Database, name string, description string, permissions []primitive.ObjectID) (models.Role, error) {
	var role models.Role

	role.ID = primitive.NewObjectID()
	role.Name = name
	role.Description = description
	role.Permissions = permissions
	role.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	// check if the role already exists
	existingRole, err := operations.GetOneByFilter[models.Role](database.Tables.Role, bson.M{"name": role.Name})
	if err == nil {
		return existingRole, nil
	}

	_, err = operations.InsertOne(database.Tables.Role, role)

	return role, err
}
