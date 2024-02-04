package models

import "gorm.io/gorm"

type Soft_skill struct {
    gorm.Model
    SkillID        uint 
    Title         string 
}
