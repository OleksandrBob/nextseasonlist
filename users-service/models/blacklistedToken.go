package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BlacklistedToken struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Token string             `bson:"token" json:"token" binding:"required"`
}
