package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPermission struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id"`
	PermissionIds []primitive.ObjectID `json:"permission_ids" bson:"permission_ids"`
}
