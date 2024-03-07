package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Role      string
	Skills    []GroupSkill `gorm:"many2many:member_skill;"`
	ProfileID uint
	Profile   Profile
	GroupID   uint
}
