package model

import (
	"fmt"

	"gorm.io/gorm"
)

type constraintTarget struct {
	model any
	name  string
}

func automigrate(db *gorm.DB) error {
	// Favorite / Comment 已是显式外键表（UserID、VideoID），不是 User/Video 上的 many2many 字段，不需要再使用 SetupJoinTable；
	if err := db.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Comment{}, &CommentFavorite{}, &Relation{}); err != nil {
		return fmt.Errorf("auto migrate failed:%w", err)
	}

	targets := []constraintTarget{
		{model: &Comment{}, name: "User"},
		{model: &Comment{}, name: "Video"},
		{model: &Comment{}, name: "ParentComment"},
		{model: &Comment{}, name: "RootComment"},
		{model: &Favorite{}, name: "User"},
		{model: &Favorite{}, name: "Video"},
		{model: &CommentFavorite{}, name: "User"},
		{model: &CommentFavorite{}, name: "Comment"},
	}
	if err := ensureCascadeConstraints(db, targets); err != nil {
		return fmt.Errorf("ensure cascade constraints failed:%w", err)
	}
	return nil
}

func ensureCascadeConstraints(db *gorm.DB, targets []constraintTarget) error {
	for _, target := range targets {
		if db.Migrator().HasConstraint(target.model, target.name) {
			if err := db.Migrator().DropConstraint(target.model, target.name); err != nil {
				return fmt.Errorf("drop constraint failed:%w", err)
			}
		}
		if err := db.Migrator().CreateConstraint(target.model, target.name); err != nil {
			return fmt.Errorf("create constraint failed:%w", err)
		}
	}
	return nil
}
