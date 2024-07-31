package models

type BookAuthor struct {
	ID        string `json:"id" bson:"_id"`
	BookId    string `json:"book_id" bson:"book_id"`
	AuthorId  string `json:"author_id" bson:"author_id"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}
