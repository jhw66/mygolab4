package model

import (
	"fmt"

	"gorm.io/gorm"
)

func automigrate(db *gorm.DB) error {
	if err := db.SetupJoinTable(&User{}, "video", &Favorite{}); err != nil {
		return fmt.Errorf("setup join table failed:%w", err)
	}
	if err := db.SetupJoinTable(&Video{}, "user", &Favorite{}); err != nil {
		return fmt.Errorf("setup join table failed:%w", err)
	}
	if err := db.SetupJoinTable(&User{}, "video", &Comment{}); err != nil {
		return fmt.Errorf("setup join table failed:%w", err)
	}
	if err := db.SetupJoinTable(&Video{}, "user", &Comment{}); err != nil {
		return fmt.Errorf("setup join table failed:%w", err)
	}
	if err := db.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Comment{}, &Relation{}); err != nil {
		return fmt.Errorf("auto migrate failed:%w", err)
	}
	return nil
}
