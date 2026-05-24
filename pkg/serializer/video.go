package serializer

import (
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
)

type Video struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	Info          string `json:"info"`
	Cover         string `json:"cover"`
	CommentCount  uint   `json:"comment_count"`
	FavoriteCount uint   `json:"favorite_count"`
	CreatedAt     int64  `json:"created_at"`
}

func VideoItemFromModel(video *model.Video) *types.VideoItem {
	return &types.VideoItem{
		Id:            video.ID,
		Title:         video.Title,
		Url:           video.URL,
		Info:          video.Info,
		Cover:         video.Cover,
		CommentCount:  video.CommentCount,
		FavoriteCount: video.FavoriteCount,
		CreatedAt:     video.CreatedAt.Unix(),
	}
}

func VideoRspFromModel(video *model.Video) *types.VideoRsp {
	item := VideoItemFromModel(video)
	return &types.VideoRsp{
		Status: 200,
		Data:   item,
	}
}

func VideoListRspFromModels(videos []model.Video) *types.VideoListRsp {
	items := make([]types.VideoItem, 0, len(videos))
	for i := range videos {
		items = append(items, *VideoItemFromModel(&videos[i]))
	}
	return &types.VideoListRsp{
		Status: 200,
		Data:   &items,
	}
}
