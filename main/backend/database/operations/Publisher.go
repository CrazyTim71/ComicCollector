package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllPublishers(db *mongo.Database) ([]models.Publisher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("publisher").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var publishers []models.Publisher
	err = cursor.All(ctx, &publishers)

	return publishers, err
}

func GetPublisherById(db *mongo.Database, id string) (models.Publisher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var publisher models.Publisher
	err := db.Collection("publisher").FindOne(ctx, bson.M{"_id": id}).Decode(&publisher)

	return publisher, err
}

func GetPublisherByName(db *mongo.Database, name string) (models.Publisher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var publisher models.Publisher
	err := db.Collection("publisher").FindOne(ctx, bson.M{"name": name}).Decode(&publisher)

	return publisher, err
}

func InsertPublisher(db *mongo.Database, publisher models.Publisher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("publisher").InsertOne(ctx, publisher)

	return err
}

func UpdatePublisher(db *mongo.Database, id string, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("publisher").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeletePublisher(db *mongo.Database, id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("publisher").DeleteOne(ctx, bson.M{"_id": id})

	return result, err
}
