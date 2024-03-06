package models

import (
	"gorm.io/gorm"
)

type GroupSkill struct {
	gorm.Model
	Name      string        `json:"name" binding:"required"`
	Positions []ReqPosition `gorm:"many2many:position_skill;"`
	Members   []Member      `gorm:"many2many:member_skill;"`
}
