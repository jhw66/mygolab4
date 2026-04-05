package model

type Relation struct {
	BaseModel
	UserID       string `gorm:"not null;index;type:varchar(32)"`
	TargetUserID string `gorm:"not null;index;type:varchar(32)"`
}

func (Relation) TableName() string {
	return "relation"
}
