package model

import (
	"fmt"

	"gorm.io/gorm"
)

func automigrate(db *gorm.DB) error {
	// Favorite / Comment 已是显式外键表（UserID、VideoID），不是 User/Video 上的 many2many 字段，不需要再使用 SetupJoinTable；
	if err := db.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Comment{}, &Relation{}); err != nil {
		return fmt.Errorf("auto migrate failed:%w", err)
	}
	return nil
}
