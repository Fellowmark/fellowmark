package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  string
}

func CreateUser(firstName string, lastName string, email string, password string) {
}
