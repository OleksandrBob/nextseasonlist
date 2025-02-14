package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserUpdateInfo struct {
	ID        primitive.ObjectID `json:"id"`
	FirstName string             `json:"firstName" binding:"required,min=1"`
	LastName  string             `json:"lastName" binding:"required,min=1"`
}
