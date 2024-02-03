package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     string
	Groups       []*Group      `gorm:"many2many:group_member;"`
	Applications []ReqPosition `gorm:"many2many:applicants;"`
	Profile      Profile
}
