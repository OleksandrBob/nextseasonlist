package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddEpisodeCommand struct {
	Name        string             `json:"name" binding:"required"`
	Season      int32              `json:"season" binding:"required,min=1"`
	Number      int32              `json:"number" binding:"required,min=1"`
	SerialId    primitive.ObjectID `json:"serialId" binding:"required"`
	ReleaseDate time.Time          `json:"releaseDate"`
}
