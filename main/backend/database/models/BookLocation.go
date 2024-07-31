package models

type BookLocation struct {
	ID         string `json:"id" bson:"_id"`
	BookId     string `json:"book_id" bson:"book_id"`
	LocationId string `json:"location_id" bson:"location_id"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
	UpdatedAt  string `json:"updated_at" bson:"updated_at"`
}
