package models

import "gorm.io/gorm"

type Programming struct {
	gorm.Model
	Tec_skillsID uint
	Programtype string
	Type    string
}