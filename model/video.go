package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserID        uint
	Title         string
	URL           string
	Info          string
	Cover         string
	View          uint
	FavoriteCount uint
	User          User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
