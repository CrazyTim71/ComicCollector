package models

type Publisher struct {
	ID         int    `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	WebsiteURL string `json:"website_url" bson:"website_url"`
	Country    string `json:"country" bson:"country"`
}
