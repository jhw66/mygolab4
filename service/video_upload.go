package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"gorm.io/gorm"
)

func UploadVideo(tx *gorm.DB, video *model.Video) (*model.Video, *serializer.Response) {
	if err := tx.Save(video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "上传视频失败",
		}
	}
	return video, nil
}
