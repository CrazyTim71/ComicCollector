package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllBooks(db *mongo.Database) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("book").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var books []models.Book
	err = cursor.All(ctx, &books)

	return books, err
}

func GetBookById(db *mongo.Database, id primitive.ObjectID) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book models.Book
	err := db.Collection("book").FindOne(ctx, bson.M{"_id": id}).Decode(&book)

	return book, err
}

func CreateBook(db *mongo.Database, book models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("book").InsertOne(ctx, book)

	return err
}
