package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllLocations(db *mongo.Database) ([]models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("location").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var locations []models.Location
	err = cursor.All(ctx, &locations)

	return locations, err
}

func GetLocationById(db *mongo.Database, id string) (models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var location models.Location
	err := db.Collection("location").FindOne(ctx, bson.M{"_id": id}).Decode(&location)

	return location, err
}

func GetLocationByName(db *mongo.Database, name string) (models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var location models.Location
	err := db.Collection("location").FindOne(ctx, bson.M{"name": name}).Decode(&location)

	return location, err
}

func InsertLocation(db *mongo.Database, location models.Location) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("location").InsertOne(ctx, location)

	return err
}

func UpdateLocation(db *mongo.Database, id string, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("location").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeleteLocation(db *mongo.Database, id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("location").DeleteOne(ctx, bson.M{"_id": id})

	return result, err
}
