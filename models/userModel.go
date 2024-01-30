package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fname string
	Lname string
	Bio string
	Tag string
	Groups []*Group `gorm:"many2many:group_member;"`
}