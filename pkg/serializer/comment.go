package serializer

import (
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
)

type Comment struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	VideoID       string `json:"video_id"`
	CommentID     string `json:"comment_id,omitempty"`
	RootID        string `json:"root_id,omitempty"`
	Content       string `json:"content"`
	FavoriteCount uint   `json:"favorite_count"`
	CreatedAt     int64  `json:"created_at"`
	User          User   `json:"user"`
}

type CommentList struct {
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Comments []Comment `json:"comments"`
}

func CommentListRspFromModels(comments []model.Comment, total int64, page int, pageSize int) *types.CommentListRsp {
	items := make([]types.CommentItem, 0, len(comments))
	for i := range comments {
		items = append(items, types.CommentItem{
			Id:            comments[i].ID,
			UserId:        comments[i].UserID,
			VideoId:       comments[i].VideoID,
			CommentId:     stringValue(comments[i].CommentID),
			RootId:        stringValue(comments[i].RootID),
			Content:       comments[i].Content,
			FavoriteCount: comments[i].FavoriteCount,
			CreatedAt:     comments[i].CreatedAt.Unix(),
			User:          UserItemFromModel(&comments[i].User),
		})
	}
	return &types.CommentListRsp{
		Status: 200,
		Data: &types.CommentListData{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Comments: items,
		},
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
