package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Gname string 
	Owner_id int 
	Topic string
	Description string
	IsHidden bool
	Limit_mem int
	Cat_id int
	Users []*User `gorm:"many2many:group_member;"`
}