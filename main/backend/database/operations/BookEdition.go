package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllBookEditions(db *mongo.Database) ([]models.BookEdition, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("book_edition").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var bookEditions []models.BookEdition
	err = cursor.All(ctx, &bookEditions)

	return bookEditions, err
}

func GetBookEditionById(db *mongo.Database, id primitive.ObjectID) (models.BookEdition, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookEdition models.BookEdition
	err := db.Collection("book_edition").FindOne(ctx, bson.M{"_id": id}).Decode(&bookEdition)

	return bookEdition, err
}

func GetBookEditionByName(db *mongo.Database, name string) (models.BookEdition, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookEdition models.BookEdition
	err := db.Collection("book_edition").FindOne(ctx, bson.M{"name": name}).Decode(&bookEdition)

	return bookEdition, err
}

func CreateBookEdition(db *mongo.Database, bookEdition models.BookEdition) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("book_edition").InsertOne(ctx, bookEdition)

	return err
}

func UpdateBookEdition(db *mongo.Database, id primitive.ObjectID, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("book_edition").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeleteBookEdition(db *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("book_edition").DeleteOne(ctx, bson.M{"_id": id})

	return err
}
