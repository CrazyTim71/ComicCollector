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

func SaveRole(db *mongo.Database, newRole models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("role").InsertOne(ctx, newRole, options.InsertOne())

	return err
}

func GetRoleById(db *mongo.Database, roleId primitive.ObjectID) (models.Role, error) {
	var role models.Role
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Collection("role").FindOne(ctx, bson.M{"_id": roleId}).Decode(&role)

	return role, err
}

func CreateRole(name string, description string) (models.Role, error) {
	var role models.Role

	role.ID = primitive.NewObjectID()
	role.Name = name
	role.Description = description
	role.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	role.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err := SaveRole(database.MongoDB, role)
	if err != nil {
		return role, err
	}

	return role, nil
}
