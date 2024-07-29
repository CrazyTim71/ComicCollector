package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var MongoDB *mongo.Database

func InitDatabase() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Println("You must set your 'MONGODB_URI' environmental variable.")
		return errors.New("you must set your 'MONGODB_URI' environmental variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
		return err
	}

	dbName := os.Getenv("MONGODB_DBNAME")
	if dbName == "" {
		log.Println("The 'MONGODB_DBNAME' environmental variable is not set. Defaulting to 'TimeCraft'.")
		dbName = "ComicCollector"
	}
	MongoDB = client.Database(dbName)
	return nil
}
