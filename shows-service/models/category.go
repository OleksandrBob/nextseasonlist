package models

type Category struct {
	ID   int    `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name" binding:"required"`
}
