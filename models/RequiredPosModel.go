package models

import (
	"gorm.io/gorm"
)

type ReqPosition struct {
	gorm.Model
	Role    string `json:"role"`
	GroupID uint
}
