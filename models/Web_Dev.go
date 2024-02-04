package models

import "gorm.io/gorm"

type WebDev struct {
	gorm.Model
	Tec_skillsID uint
	WebDev string
	
}