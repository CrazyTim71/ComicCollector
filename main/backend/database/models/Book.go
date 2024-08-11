package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Number      int                `json:"number" bson:"number"`
	ReleaseDate string             `json:"release_date" bson:"release_date"`
	CoverImage  []byte             `json:"cover_image" bson:"cover_image"`
	Description string             `json:"description" bson:"description"`
	Notes       string             `json:"notes" bson:"notes"`
	BookType    primitive.ObjectID `json:"book_type" bson:"book_type"`
	BookEdition primitive.ObjectID `json:"book_edition" bson:"book_edition"`
	Printing    string             `json:"printing" bson:"printing"`
	ISBN        string             `json:"isbn" bson:"isbn"`
	Price       string             `json:"price" bson:"price"`
	Count       int                `json:"count" bson:"count"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
