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

func GetAllRoles(db *mongo.Database) ([]models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("role").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	err = cursor.All(ctx, &roles)

	return roles, err
}

func GetRoleById(db *mongo.Database, roleId primitive.ObjectID) (models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var role models.Role
	err := db.Collection("role").FindOne(ctx, bson.M{"_id": roleId}).Decode(&role)

	return role, err
}

func GetRoleByName(db *mongo.Database, name string) (models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var role models.Role
	err := db.Collection("role").FindOne(ctx, bson.M{"name": name}).Decode(&role)

	return role, err
}

func CreateRole(db *mongo.Database, name string, description string, permissions []primitive.ObjectID) (models.Role, error) {
	var role models.Role

	role.ID = primitive.NewObjectID()
	role.Name = name
	role.Description = description
	role.Permissions = permissions
	role.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	role.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	// check if the role already exists
	existingRole, err := GetRoleByName(db, name)
	if err == nil {
		return existingRole, nil
	}

	err = InsertRole(database.MongoDB, role)

	return role, err
}

func CheckIfAllRolesExist(db *mongo.Database, roleIds []primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("role").Find(ctx, bson.M{"_id": bson.M{"$in": roleIds}})
	if err != nil {
		return false
	}

	var foundRoles []models.Role
	err = cursor.All(ctx, &foundRoles)

	return len(foundRoles) == len(roleIds)
}

func InsertRole(db *mongo.Database, newRole models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("role").InsertOne(ctx, newRole, options.InsertOne())

	return err
}

func UpdateRole(db *mongo.Database, roleId primitive.ObjectID, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("role").UpdateOne(ctx, bson.M{"_id": roleId}, bson.M{"$set": data})

	return result, err
}

func DeleteRole(db *mongo.Database, roleId primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("role").DeleteOne(ctx, bson.M{"_id": roleId})

	return result, err
}
