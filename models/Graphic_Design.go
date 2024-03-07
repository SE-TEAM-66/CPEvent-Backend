package models

import "gorm.io/gorm"

type GraphicDesign struct {
	gorm.Model
	SkillID uint
	GraphicDesign string
	
}