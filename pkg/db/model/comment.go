package model

type Comment struct {
	BaseModel
	UserID  string `gorm:"index;not null;type:varchar(32)"`
	VideoID string `gorm:"index;not null;type:varchar(32)"`
	Content string

	User  User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Video Video `gorm:"foreignKey:VideoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Comment) TableName() string {
	return "comment"
}
