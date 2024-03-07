package models

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
	ProfileID  uint
	Profile Profile `gorm:"foreignKey:ProfileID"`
	Soft_skill []Soft_skill 
	Lang_skill []Lang_skill
	DataAna []DataAna
	DBmanage []DBmanage
	GraphicDesign []GraphicDesign
	Programming []Programming
	WebDev []WebDev
	Type string
}
