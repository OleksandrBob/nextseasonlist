package models

type RegisterUserDto struct {
	FirstName string `bson:"firstName" json:"firstName" binding:"required,min=1"`
	LastName  string `bson:"lastName" json:"lastName" binding:"required,min=1"`
	Email     string `bson:"email" json:"email" binding:"required,email"`
	Password  string `bson:"password" json:"password" binding:"required,min=6"`
}
