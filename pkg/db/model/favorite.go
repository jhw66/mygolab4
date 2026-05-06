package model

type Favorite struct {
	UserID  string `gorm:"primaryKey;not null;type:varchar(32)"`
	VideoID string `gorm:"primaryKey;not null;type:varchar(32)"`

	User  User  `gorm:"foreignKey:UserID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Video Video `gorm:"foreignKey:VideoID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Favorite) TableName() string {
	return "favorite"
}
