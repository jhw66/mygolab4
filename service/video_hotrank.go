package service

import "github.com/jhw66/myvideo_lab4/model"

func GetRankVideos(limit int) ([]model.Video, error) {
	var videos []model.Video
	err := model.Db.Order("favorite_count desc").Order("view desc").Order("created_at desc").
		Limit(limit).Find(&videos).Error
	return videos, err
}
