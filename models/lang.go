package models

import "gorm.io/gorm"

type Lang_skill struct {
	gorm.Model
	SkillID uint
	Title string
}