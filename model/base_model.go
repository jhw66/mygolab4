package model

import (
	"time"

	"github.com/jhw66/myvideo_lab4/utils"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string `gorm:"primaryKey;type:varchar(32)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == " " {
		id, err := utils.GenerateID()
		if err != nil {
			return err
		}
		b.ID = id
	}
	return nil
}
