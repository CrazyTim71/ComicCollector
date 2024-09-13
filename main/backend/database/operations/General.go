package operations

import (
	"ComicCollector/main/backend/database"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CheckIfExists(db *mongo.Database, collectionName string, filter bson.M) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(collectionName)
	_, err := collection.Find(ctx, filter)
	if err != nil {
		return false
	}
	return true
}

func GetAll[T any](collectionName string) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.MongoDB.Collection(collectionName).Find(ctx, bson.M{})
	var result []T
	err = cursor.All(ctx, &result)

	return result, err
}

func GetOneById[T any](collectionName string, id primitive.ObjectID) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := database.MongoDB.Collection(collectionName).FindOne(ctx, bson.M{"_id": id})

	var decodedResult T
	if err := result.Decode(&decodedResult); err != nil {
		return decodedResult, err
	}

	return decodedResult, nil
}

func GetOneByFilter[T any](collectionName string, filter bson.M) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := database.MongoDB.Collection(collectionName).FindOne(ctx, filter)

	var decodedResult T
	if err := result.Decode(&decodedResult); err != nil {
		return decodedResult, err
	}

	return decodedResult, nil
}

func GetManyByFilter[T any](collectionName string, filter bson.M) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.MongoDB.Collection(collectionName).Find(ctx, filter)
	var result []T
	err = cursor.All(ctx, &result)

	return result, err
}

func GetManyIdsByFilter(collectionName string, idName string, filter bson.M) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.MongoDB.Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for cursor.Next(ctx) {
		var result bson.M
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		value := result[idName]
		if value == nil {
			return nil, nil
		}

		switch v := value.(type) {
		case primitive.A: // array of ObjectIDs
			for _, item := range v {
				if objID, ok := item.(primitive.ObjectID); ok {
					ids = append(ids, objID)
				} else {
					return nil, fmt.Errorf("expected primitive.ObjectID, got %T", item)
				}
			}
		default: // single ObjectID
			ids = append(ids, result[idName].(primitive.ObjectID))
		}
	}

	return ids, nil
}

func InsertOne(collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.MongoDB.Collection(collectionName).InsertOne(ctx, data)
	return result, err
}

func UpdateOne(collectionName string, id primitive.ObjectID, data bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.MongoDB.Collection(collectionName).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})

	return result, err
}

func DeleteOne(collectionName string, filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.MongoDB.Collection(collectionName).DeleteOne(ctx, filter)

	return result, err
}

func DeleteMany(collectionName string, filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.MongoDB.Collection(collectionName).DeleteMany(ctx, filter)

	return result, err
}

func CheckIfAllIdsExist[T any](collectionName string, ids []primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.MongoDB.Collection(collectionName).Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return false
	}

	var result []T
	err = cursor.All(ctx, &result)

	return len(result) == len(ids)
}
