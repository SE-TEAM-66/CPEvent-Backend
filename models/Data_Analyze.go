package models

import "gorm.io/gorm"

type DataAna struct {
	gorm.Model
	Tec_skillsID uint
	DataAna string
}