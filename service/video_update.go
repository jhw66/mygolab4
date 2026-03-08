package service

import (
	"github.com/jhw66/myvideo_lab4/model"
)

func CompareVidAndUid(uid uint, vid uint) bool {
	var video model.Video
	model.Db.Model(&video).Where("id = ?", vid).Take(&video)
	return video.UserID == uid
}
