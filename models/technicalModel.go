package models

import "gorm.io/gorm"

type Tec_skills struct {
	gorm.Model
	SkillID uint
	DataAna []DataAna
	DBmanage []DBmanage
	GraphicDesign []GraphicDesign
	Programming []Programming
	WebDev []WebDev
	Type string
}