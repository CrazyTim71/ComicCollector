package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	Username  string               `json:"username" bson:"username"`
	Password  string               `json:"password" bson:"password"`
	Roles     []primitive.ObjectID `json:"roles" bson:"roles"`
	CreatedAt primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime   `json:"updated_at" bson:"updated_at"`
	CreatedBy primitive.ObjectID   `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID   `json:"updated_by" bson:"updated_by"`
}
