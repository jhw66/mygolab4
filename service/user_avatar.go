package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
)

func UploadAvatar(user *model.User) (*model.User, *serializer.Response) {
	if err := model.Db.Save(user).Error; err != nil {
		return user, &serializer.Response{
			Status: 500,
			Msg:    "更新用户头像失败",
		}
	}
	return user, nil
}
