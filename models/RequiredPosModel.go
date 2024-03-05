package models

import (
	"gorm.io/gorm"
)

type ReqPosition struct {
	gorm.Model
	Role       string `json:"role" binding:"required"`
	GroupID    uint
	Applicants []User       `gorm:"many2many:applicants;"`
	Skills     []GroupSkill `gorm:"many2many:applicant_skill;"`
}
