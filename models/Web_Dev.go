package models

import "gorm.io/gorm"

type WebDev struct {
	gorm.Model
	SkillID uint
	WebDev string
	
}