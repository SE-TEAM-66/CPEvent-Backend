package models

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
	ProfileID  uint
	Profile Profile `gorm:"foreignKey:ProfileID"`
	Soft_skill []Soft_skill 
	Tec_skills []Tec_skills
	Lang_skill []Lang_skill
	Type string
}
