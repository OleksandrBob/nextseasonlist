package models

type AddSerialCommand struct {
	Title       string   `json:"name" binding:"required,min=1"`
	Description string   `json:"description" binding:"required,min=1"`
	Categories  []string `json:"categories"`
}
