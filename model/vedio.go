package model

import "gorm.io/gorm"

type Vedio struct {
	gorm.Model
	UserID uint
	Title  string
	URL    string
	Info   string
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
