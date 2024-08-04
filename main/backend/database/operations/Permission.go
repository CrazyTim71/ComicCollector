package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SavePermission(db *mongo.Database, newPermission models.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("permission").InsertOne(ctx, newPermission, options.InsertOne())
	if err != nil {
		return err
	}

	return nil
}
