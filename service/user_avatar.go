package service

import (
	"os"
	"strings"

	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"gorm.io/gorm"
)

func UploadAvatar(user *model.User, oldAvatarPath string) (*model.User, *serializer.Response) {
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(user).Error
	})
	if err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "更新用户头像失败",
		}
	}

	if oldAvatarPath != "" {
		oldpath := strings.TrimPrefix(oldAvatarPath, "/")
		if err := os.Remove(oldpath); err != nil && !os.IsNotExist(err) {
			return nil, &serializer.Response{
				Status: 500,
				Msg:    "删除旧头像文件失败",
			}
		}
	}

	return user, nil
}
