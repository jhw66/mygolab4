package model

type ChatMember struct {
	UserID   string   `gorm:"PrimaryKey;not null;type:varchar(32)"`
	RoomID   string   `gorm:"PrimaryKey;not null;type:varchar(32)"`
	Role     string   `gorm:"not null;default:member"`
	ChatRoom ChatRoom `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ChatMember) TableName() string {
	return "chat_room_member"
}
