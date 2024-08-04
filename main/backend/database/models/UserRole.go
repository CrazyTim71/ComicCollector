package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRole struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	RoleId    primitive.ObjectID `json:"role_id" bson:"role_id"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
