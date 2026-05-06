package model

type Video struct {
	BaseModel
	UserID        string
	Title         string
	URL           string
	Info          string
	Cover         string
	CommentCount  uint
	FavoriteCount uint
	HotScore      uint64 `gorm:"index:idx_video_hot_rank,priority:1"`
	User          User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Video) TableName() string {
	return "video"
}
