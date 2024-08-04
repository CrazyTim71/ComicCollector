package operations

import (
	"ComicCollector/main/backend/database/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SaveRolePermission(db *mongo.Database, newRolePermission models.RolePermission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("role_permission").InsertOne(ctx, newRolePermission, options.InsertOne())
	if err != nil {
		return err
	}

	return nil
}
