package model

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	UserID       uint `gorm:"not null;index"`
	TargetUserID uint `gorm:"not null;index"`
}

func (Relation) TableName() string {
	return "relation"
}
