package model

type CommentFavorite struct {
	UserID    string `gorm:"primaryKey;not null;type:varchar(32)"`
	CommentID string `gorm:"primaryKey;not null;type:varchar(32)"`

	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Comment Comment `gorm:"foreignKey:CommentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CommentFavorite) TableName() string {
	return "comment_favorite"
}
