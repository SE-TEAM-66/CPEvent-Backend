package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	// PMaw_ID uint `gorm:"primary_key"`
	PicUrl string
	Etitle string
	Edesc  string
	Edate  string
	Etime  string
	Groups []Group `gorm:"foreignKey:EventID"`
}