package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RolePermission struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"` // TODO: remove the name ??
	RoleId       primitive.ObjectID `json:"role_id" bson:"role_id"`
	PermissionId primitive.ObjectID `json:"permission_id" bson:"permission_id"`
	CreatedAt    primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
