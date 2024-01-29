package models

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"gorm.io/gorm"
)

type Position struct {
	gorm.Model
	Role    string `json:"role"`
	GroupID uint
}

func AddPosition(gid uint, role string) error {
	newPos := Position{Role: role, GroupID: gid}
	result := initializers.DB.Create(&newPos)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeletePosition(pid uint) error {
	var position Position
	result := initializers.DB.Delete(&position, pid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
