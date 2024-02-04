package models

import "gorm.io/gorm"

type GraphicDesign struct {
	gorm.Model
	Tec_skillsID uint
	GraphicDesign string
	
}