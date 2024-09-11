package helpers

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils"
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
	existingPermission, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
	if err == nil {
		return existingPermission, nil
	}

	err = operations.InsertPermission(database.MongoDB, permission)
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
	existingRole, err := operations.GetRoleByName(db, name)
	if err == nil {
		return existingRole, nil
	}

	err = operations.InsertRole(database.MongoDB, role)

	return role, err
}
