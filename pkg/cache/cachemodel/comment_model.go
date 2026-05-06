package cachemodel

import "github.com/jhw66/myvideo_lab4/pkg/db/model"

type CommentListCacheData struct {
	Comments []model.Comment `json:"comments"`
	Total    int64           `json:"total"`
}
