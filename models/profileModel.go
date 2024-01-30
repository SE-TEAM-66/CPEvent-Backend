package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model 
	ProfilePicture string
	Fname string
	Lname string
	Faculty string
	Bio string
	Phone string
	Email string
	Facebook string
	Line string
	UserID uint
}