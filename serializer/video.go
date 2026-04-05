package serializer

import "github.com/jhw66/myvideo_lab4/model"

type Video struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	Info          string `json:"info"`
	Cover         string `json:"avatar"`
	CommentCount  uint   `json:"comment_counts"`
	FavoriteCount uint   `json:"favorite_counts"`
	CreatedAt     int64  `json:"created_at"`
}

func BuildVideo(video *model.Video) *Video {
	return &Video{
		ID:            video.ID,
		Title:         video.Title,
		URL:           video.URL,
		Info:          video.Info,
		Cover:         video.Cover,
		CommentCount:  video.CommentCount,
		FavoriteCount: video.FavoriteCount,
		CreatedAt:     video.CreatedAt.Unix(),
	}
}

func BuildVideoList(videos *[]model.Video) *[]Video {
	var res []Video
	for key := range *videos {
		res = append(res, *BuildVideo(&(*videos)[key]))
	}
	return &res
}
