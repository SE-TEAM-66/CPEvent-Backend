package models

import (
	"gorm.io/gorm"
)

type ReqPosition struct {
	gorm.Model
	Role       string `json:"role"`
	GroupID    uint
	Applicants []User `gorm:"many2many:applicants;"`
}
