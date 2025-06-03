package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SearchSerialsQuery struct {
	Title      string   `bson:"name" json:"name"`
	Categories []string `bson:"categories" json:"categories"`
}

type SearchSerialsQueryResponse struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string             `bson:"title" json:"title"`
}
