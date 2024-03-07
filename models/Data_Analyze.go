package models

import "gorm.io/gorm"

type DataAna struct {
	gorm.Model
	SkillID uint
	DataAna string
	
}
