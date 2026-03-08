package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
)

func FindVideoByUser(user *model.User) (*[]model.Video, *serializer.Response) {
	var videos []model.Video
	if err := model.Db.Model(user).Association("Videos").Find(&videos); err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "查找失败",
		}
	}
	return &videos, nil
}
