package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"gorm.io/gorm"
)

func DeleteVideo(tx *gorm.DB, vid uint) (*model.Video, *serializer.Response) {
	var video model.Video
	tx.Where("id = ?", vid).Take(&video)
	if err := tx.Delete(&video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "删除视频失败",
		}
	}
	return &video, nil
}
