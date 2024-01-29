package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Gname string `json`
	Owner_id int
	Topic string
	Description string
	IsHidden bool
	Limit_mem int
	Cat_id int
}