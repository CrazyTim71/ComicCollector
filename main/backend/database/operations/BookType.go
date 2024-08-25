package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllBookTypes(db *mongo.Database) ([]models.BookType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("book_type").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var bookTypes []models.BookType
	err = cursor.All(ctx, &bookTypes)

	return bookTypes, err
}

func GetBookTypeById(db *mongo.Database, id primitive.ObjectID) (models.BookType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookType models.BookType
	err := db.Collection("book_type").FindOne(ctx, bson.M{"_id": id}).Decode(&bookType)

	return bookType, err
}

func GetBookTypeByName(db *mongo.Database, name string) (models.BookType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookType models.BookType
	err := db.Collection("book_type").FindOne(ctx, bson.M{"name": name}).Decode(&bookType)

	return bookType, err
}

func InsertBookType(db *mongo.Database, bookType models.BookType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("book_type").InsertOne(ctx, bookType)

	return err
}

func UpdateBookType(db *mongo.Database, id primitive.ObjectID, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("book_type").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeleteBookType(db *mongo.Database, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("book_type").DeleteOne(ctx, bson.M{"_id": id})

	return result, err
}
