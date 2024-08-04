package database

import (
	"ComicCollector/main/backend/utils/env"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var MongoDB *mongo.Database

func InitDatabase() bool {
	uri := env.GetDatabaseURI()
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	dbName := env.GetDatabaseName()
	if dbName == "" {
		log.Fatal("The 'MONGODB_DBNAME' environmental variable is not set.")
	}
	MongoDB = client.Database(dbName)

	return true
}

func HasCollection(db *mongo.Database, collectionName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.D{{"name", collectionName}})
	if err != nil {
		log.Println(err)
		return false
	}

	return len(collections) == 1
}
