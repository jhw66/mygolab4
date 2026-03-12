package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint `gorm:"index;not null"`
	VideoID uint `gorm:"index;not null"`
	Content string

	User  User  `gorm:"foreignKey:UserID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Video Video `gorm:"foreignKey:VideoID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Comment) TableName() string {
	return "comment"
}
