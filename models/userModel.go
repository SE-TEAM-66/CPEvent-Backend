package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fname    string
	Lname    string
	Email    string `gorm:"unique"`
	Password string
	Bio      string
	Tag      string
	Profile Profile `gorm:"foreignKey:UserID"`
}