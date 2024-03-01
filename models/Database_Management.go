package models

import "gorm.io/gorm"

type DBmanage struct {
	gorm.Model
	SkillID uint
	DBmanage string
	
}