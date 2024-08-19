package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllAuthors(db *mongo.Database) ([]models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("author").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var authors []models.Author
	err = cursor.All(ctx, &authors)

	return authors, err
}

func GetAuthorById(db *mongo.Database, id primitive.ObjectID) (models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var author models.Author
	err := db.Collection("author").FindOne(ctx, bson.M{"_id": id}).Decode(&author)

	return author, err
}

func GetAuthorByName(db *mongo.Database, name string) (models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var author models.Author
	err := db.Collection("author").FindOne(ctx, bson.M{"name": name}).Decode(&author)

	return author, err
}

func CreateAuthor(db *mongo.Database, author models.Author) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("author").InsertOne(ctx, author)

	return err
}

func UpdateAuthor(db *mongo.Database, id primitive.ObjectID, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("author").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeleteAuthor(db *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("author").DeleteOne(ctx, bson.M{"_id": id})

	return err
}
