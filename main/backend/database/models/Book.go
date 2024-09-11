package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Title       string               `json:"title" bson:"title"`
	Number      int                  `json:"number" bson:"number"`
	ReleaseDate primitive.DateTime   `json:"release_date" bson:"release_date"`
	CoverImage  primitive.ObjectID   `json:"cover_image" bson:"cover_image"` // GridFS Id
	Description string               `json:"description" bson:"description"`
	Notes       string               `json:"notes" bson:"notes"`
	Authors     []primitive.ObjectID `json:"authors" bson:"authors"`
	Publishers  []primitive.ObjectID `json:"publishers" bson:"publishers"`
	Locations   []primitive.ObjectID `json:"locations" bson:"locations"`
	Owners      []primitive.ObjectID `json:"owners" bson:"owners"`
	BookType    primitive.ObjectID   `json:"book_type" bson:"book_type"`
	BookEdition primitive.ObjectID   `json:"book_edition" bson:"book_edition"`
	Printing    string               `json:"printing" bson:"printing"`
	ISBN        string               `json:"isbn" bson:"isbn"`
	Price       string               `json:"price" bson:"price"`
	Count       int                  `json:"count" bson:"count"`
	CreatedAt   primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime   `json:"updated_at" bson:"updated_at"`
	CreatedBy   primitive.ObjectID   `json:"created_by" bson:"created_by"`
	UpdatedBy   primitive.ObjectID   `json:"updated_by" bson:"updated_by"`
}
