package models

type BookPublisher struct {
	ID          string `json:"id" bson:"_id"`
	BookId      string `json:"book_id" bson:"book_id"`
	PublisherId string `json:"publisher_id" bson:"publisher_id"`
	CreatedAt   string `json:"created_at" bson:"created_at"`
	UpdatedAt   string `json:"updated_at" bson:"updated_at"`
}
