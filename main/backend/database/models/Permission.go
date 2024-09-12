package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy   primitive.ObjectID `json:"updated_by" bson:"updated_by"`
}
