package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	PicUrl string
	Etitle string
	Edesc  string
	Edate  string
	Etime  string
	Eloc   string
	Groups []Group `gorm:"foreignKey:EventID"`
}