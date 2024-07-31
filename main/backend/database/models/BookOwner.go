package models

type BookOwner struct {
	ID        string `json:"id" bson:"_id"`
	BookId    string `json:"book_id" bson:"book_id"`
	OwnerId   string `json:"owner_id" bson:"owner_id"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}
