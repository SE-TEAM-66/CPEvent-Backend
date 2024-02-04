package models

import "gorm.io/gorm"

type DBmanage struct {
	gorm.Model
	Tec_skillsID uint
	DBmanage string
}