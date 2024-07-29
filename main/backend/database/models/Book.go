package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Title       string               `json:"title" bson:"title"`
	ReleaseDate string               `json:"release_date" bson:"release_date"`
	AddedDate   string               `json:"added_date" bson:"added_Date"`
	CoverImage  string               `json:"cover_image" bson:"cover_image"`
	Description string               `json:"description" bson:"description"`
	Notes       string               `json:"notes" bson:"notes"`
	Publishers  []primitive.ObjectID `json:"publisher" bson:"publisher"`
	Authors     []primitive.ObjectID `json:"authors" bson:"authors"`
	BookType    primitive.ObjectID   `json:"book_type" bson:"book_type"`
	BookEdition primitive.ObjectID   `json:"book_edition" bson:"book_edition"`
	Printing    string               `json:"printing" bson:"printing"`
	ISBN        string               `json:"isbn" bson:"isbn"`
	Price       string               `json:"price" bson:"price"`
	Count       int                  `json:"count" bson:"count"`
	Location    primitive.ObjectID   `json:"location" bson:"location"`
	Owner       primitive.ObjectID   `json:"owner" bson:"owner"`
}
