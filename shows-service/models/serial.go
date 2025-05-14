package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Serial struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Description string             `bson:"description" json:"description"`
	Categories  []string           `bson:"categories" json:"categories"`
}
