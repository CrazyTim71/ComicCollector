package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RolePermission struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}
