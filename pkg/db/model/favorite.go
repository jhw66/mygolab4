package model

type Favorite struct {
	UserID  string `gorm:"primaryKey;not null;type:varchar(32)"`
	VideoID string `gorm:"primaryKey;not null;type:varchar(32)"`

	User  User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Video Video `gorm:"foreignKey:VideoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Favorite) TableName() string {
	return "favorite"
}
