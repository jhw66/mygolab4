package model

type Comment struct {
	BaseModel
	UserID        string  `gorm:"index;not null;type:varchar(32)"`
	VideoID       string  `gorm:"index;not null;type:varchar(32)"`
	CommentID     *string `gorm:"index;type:varchar(32)"`
	RootID        *string `gorm:"index;type:varchar(32)"`
	Content       string
	FavoriteCount uint

	User          User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Video         Video    `gorm:"foreignKey:VideoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment *Comment `gorm:"foreignKey:CommentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RootComment   *Comment `gorm:"foreignKey:RootID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Comment) TableName() string {
	return "comment"
}
