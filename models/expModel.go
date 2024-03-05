package models

import "gorm.io/gorm"

type Exp struct {
	gorm.Model
	ProfileID   uint
	Profile     Profile `gorm:"foreignKey:ProfileID"`
	Title       string
	Description string
}
