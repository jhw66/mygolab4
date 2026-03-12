package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"gorm.io/gorm"
)

type VideoSearch struct {
	KeyWord string `form:"key_word" json:"key_word" binding:"required,max=10"`
}

func (videosearch VideoSearch) FindVideosByKeyword() (*[]model.Video, *serializer.Response) {
	var videos []model.Video
	var count int64
	if err := model.Db.Where("title like ?", "%"+videosearch.KeyWord+"%").Find(&videos).
		Count(&count).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "查找失败",
		}
	}
	if count == 0 {
		return nil, &serializer.Response{
			Status: 404,
			Msg:    "未找到相关视频",
		}
	}
	return &videos, nil
}

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

func GetRankVideos(limit int) ([]model.Video, error) {
	var videos []model.Video
	err := model.Db.Order("favorite_count desc").Order("comment_count desc").Order("created_at desc").
		Limit(limit).Find(&videos).Error
	return videos, err
}

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

func FindVideoByVid(vid uint) (*model.Video, *serializer.Response) {
	var video model.Video
	if err := model.Db.Where("id = ?", vid).Take(&video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 404,
			Msg:    "未找到该视频",
		}
	}
	return &video, nil
}

func CompareVidAndUid(uid uint, vid uint) bool {
	var video model.Video
	model.Db.Model(&video).Where("id = ?", vid).Take(&video)
	return video.UserID == uid
}

func UploadVideo(tx *gorm.DB, video *model.Video) (*model.Video, *serializer.Response) {
	if err := tx.Save(video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "上传视频失败",
		}
	}
	return video, nil
}
