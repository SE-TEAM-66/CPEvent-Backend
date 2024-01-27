package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fname string
	Lname string
	Bio string
	Tag string
}