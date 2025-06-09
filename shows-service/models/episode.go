package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	SerialId    primitive.ObjectID `bson:"serialId" json:"serialId" binding:"required"`
	Season      int32              `bson:"season" json:"season" binding:"required"`
	Number      int32              `bson:"number" json:"number" binding:"required"`
	ReleaseDate time.Time          `bson:"releaseDate" json:"releaseDate"`
}
