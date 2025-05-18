package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName" binding:"required,min=1"`
	LastName  string             `bson:"lastName" json:"lastName" binding:"required,min=1"`
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	Password  string             `bson:"password" json:"password" binding:"required,min=6"`
	Roles     []string           `bson:"roles" json:"roles" binding:"required"`
}
