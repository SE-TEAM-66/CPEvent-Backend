package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     string
	Applications []ReqPosition `gorm:"many2many:applicants;"`
	Profile      Profile
}
