package operations

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllPermissionsFromRole(db *mongo.Database, roleId primitive.ObjectID) ([]models.Permission, error) {
	var permissions []models.Permission

	role, err := GetOneById[models.Role](database.Tables.Role, roleId)
	if err != nil {
		return permissions, err
	}

	permissions, err = GetManyByFilter[models.Permission](database.Tables.Permission, bson.M{"_id": bson.M{"$in": role.Permissions}})

	return permissions, err
}
