package models

type SearchSerialsQuery struct {
	Name       string   `bson:"name" json:"name"`
	Categories []string `bson:"categories" json:"categories"`
}
