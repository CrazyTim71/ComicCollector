package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Publisher struct {
	ID         int                `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	WebsiteURL string             `json:"website_url" bson:"website_url"`
	Country    string             `json:"country" bson:"country"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt  primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
