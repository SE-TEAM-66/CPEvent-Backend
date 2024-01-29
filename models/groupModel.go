package models

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Gname       string
	Owner_id    int
	Topic       string
	Description string
	IsHidden    bool
	Limit_mem   int
	Cat_id      int
	Positions   []Position `gorm:"foreignKey:GroupID"`
}

func GetPosition(gid uint) ([]Position, error) {
	var positions []Position
	result := initializers.DB.Where("group_id = ?", gid).Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
}

func ValidPosition(gid uint, pid uint) error {
	var positions []Position
	result := initializers.DB.Where("group_id = ?", gid).First(&positions, pid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
