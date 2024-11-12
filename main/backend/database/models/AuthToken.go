package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// This enum maps the type of the token
const (
	AuthRefreshToken = iota
	AuthAPIToken
)

type AuthToken struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Type      int                `json:"type" bson:"type"`
	Token     string             `json:"token" bson:"token"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	ExpiresAt int64              `json:"expiresAt" bson:"expiresAt"`
	Expires   bool               `json:"expires" bson:"expires"`
	IssuedAt  int64              `json:"issuedAt" bson:"issuedAt"`
}
