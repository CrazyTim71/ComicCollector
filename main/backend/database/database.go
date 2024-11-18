package database

import (
	"ComicCollector/main/backend/utils/env"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var MongoDB *mongo.Database
var CoverBucket *gridfs.Bucket

type TableStruct struct {
	Author      string
	Book        string
	BookEdition string
	BookType    string
	Image       string
	Location    string
	Owner       string
	Permission  string
	Publisher   string
	Role        string
	User        string
	AuthToken   string
}

var Tables = TableStruct{
	Author:      "author",
	Book:        "book",
	BookEdition: "book_edition",
	BookType:    "book_type",
	Image:       "image",
	Location:    "location",
	Owner:       "owner",
	Permission:  "permission",
	Publisher:   "publisher",
	Role:        "role",
	User:        "user",
	AuthToken:   "auth_token",
}

func InitDatabase() bool {
	uri := env.GetDatabaseURI()
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	dbName := env.GetDatabaseName()
	if dbName == "" {
		log.Fatal("The 'MONGODB_DBNAME' environmental variable is not set.")
	}
	MongoDB = client.Database(dbName)

	opt := options.GridFSBucket()
	opt.SetName("covers")
	CoverBucket, err = gridfs.NewBucket(MongoDB, opt)
	if err != nil {
		log.Fatal(err)
	}

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
