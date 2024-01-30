package models

import (
	"gorm.io/gorm"
)

type Position struct {
	gorm.Model
	Role    string `json:"role"`
	GroupID uint
}
