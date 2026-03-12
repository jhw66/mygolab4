package model

type Favorite struct {
	UserID  uint `gorm:"PrimaryKey;not null"`
	VideoID uint `gorm:"PrimaryKey;not null"`

	User  User  `gorm:"foreignKey:UserID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Video Video `gorm:"foreignKey:VideoID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Favorite) TableName() string {
	return "favorite"
}
