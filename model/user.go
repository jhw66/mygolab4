package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName       string `gorm:"unique;not null"`
	NickName       string `gorm:"not null"`
	PasswordDigest string
	Avatar         string
	Vedios         []Vedio `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func GetUser(id interface{}) User {
	var user User
	Db.Where("id = ?", id).Find(&user)
	return user
}
