package models

type AddSerialCommand struct {
	Name        string   `bson:"name" json:"name" binding:"required,min=1"`
	Description string   `bson:"description" json:"description" binding:"required,min=1"`
	Categories  []string `bson:"categories" json:"categories"`
}
