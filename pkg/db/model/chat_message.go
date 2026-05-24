package model

type ChatMessage struct {
	BaseModel
	RoomID   string   `gorm:"index;not null;type:varchar(32)"`
	UserID   string   `gorm:"index;not null;type:varchar(32)"`
	Content  string   `gorm:"not null;type:text"`
	ChatRoom ChatRoom `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
