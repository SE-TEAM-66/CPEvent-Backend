package models

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Gname        string
	Owner_id     int
	Topic        string
	Description  string
	IsHidden     bool
	Limit_mem    int
	Cat_id       int
	ReqPositions []ReqPosition `gorm:"foreignKey:GroupID"`
}
