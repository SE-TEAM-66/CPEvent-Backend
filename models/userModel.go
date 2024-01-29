package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fname string
	Lname string
	Faculty string
	Bio string
	Phone string
	Password string
}