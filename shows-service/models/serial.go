package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Serial struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Description string             `bson:"description" json:"description"`
	Categories  []string           `bson:"categories" json:"categories"`
	Seasons     int32              `bson:"seasons" json:"seasons"`
}
