package models

import (
	"gorm.io/gorm"
)

type Notify struct {
	gorm.Model
	Rec_id uint
	Sender string
	Message string//accept or not
	IsRead    bool// read or not
}
