package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRole struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	RoleId primitive.ObjectID `json:"role_id" bson:"role_id"`
}
