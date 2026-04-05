package model

type Comment struct {
	BaseModel
	UserID  string `gorm:"index;not null;type:varchar(32)"`
	VideoID string `gorm:"index;not null;type:varchar(32)"`
	Content string

	User  User  `gorm:"foreignKey:UserID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Video Video `gorm:"foreignKey:VideoID;references:ID;constrain:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Comment) TableName() string {
	return "comment"
}
