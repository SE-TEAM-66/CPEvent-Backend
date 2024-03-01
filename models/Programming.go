package models

import "gorm.io/gorm"

type Programming struct {
	gorm.Model
	SkillID uint
	Programtype string
	
}