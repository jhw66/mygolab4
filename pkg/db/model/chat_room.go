package model

type ChatRoom struct {
	BaseModel
	RoomName     string        `gorm:"not null;type:varchar(128)"`
	Visibility   string        `gorm:"not null;default:public"`
	PasswordHash string        `gorm:"type:varchar(255)"`
	OwnerID      string        `gorm:"not null;type:varchar(32)"`
	ChatMembers  []ChatMember  `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChatMessages []ChatMessage `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ChatRoom) TableName() string {
	return "chat_room"
}
