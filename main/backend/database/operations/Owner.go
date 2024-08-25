package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllOwners(db *mongo.Database) ([]models.Owner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("owner").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var owners []models.Owner
	err = cursor.All(ctx, &owners)

	return owners, err
}

func GetOwnerById(db *mongo.Database, id string) (models.Owner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var owner models.Owner
	err := db.Collection("owner").FindOne(ctx, bson.M{"_id": id}).Decode(&owner)

	return owner, err
}

func GetOwnerByName(db *mongo.Database, name string) (models.Owner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var owner models.Owner
	err := db.Collection("owner").FindOne(ctx, bson.M{"name": name}).Decode(&owner)

	return owner, err
}

func InsertOwner(db *mongo.Database, owner models.Owner) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("owner").InsertOne(ctx, owner)

	return err
}

func UpdateOwner(db *mongo.Database, id string, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("owner").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeleteOwner(db *mongo.Database, id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("owner").DeleteOne(ctx, bson.M{"_id": id})

	return result, err
}
