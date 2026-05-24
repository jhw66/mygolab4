package serializer

import (
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
)

func CommentItemFromModel(comment *model.Comment, content string, favoriteCount uint) *types.CommentItem {
	return &types.CommentItem{
		Id:            comment.ID,
		UserId:        comment.UserID,
		VideoId:       comment.VideoID,
		CommentId:     stringValue(comment.CommentID),
		RootId:        stringValue(comment.RootID),
		Content:       content,
		FavoriteCount: favoriteCount,
		IsDeleted:     comment.DeletedAt.Valid,
		CreatedAt:     comment.CreatedAt.Unix(),
		User:          UserItemFromModel(&comment.User),
	}
}

func CommentListRspFromModels(comments []model.Comment, total int64, page int, pageSize int) *types.CommentListRsp {
	items := make([]types.CommentItem, 0, len(comments))
	for i := range comments {
		content := comments[i].Content
		favoriteCount := comments[i].FavoriteCount
		isDeleted := comments[i].DeletedAt.Valid
		if isDeleted {
			content = "该评论已删除"
			favoriteCount = 0
		}
		items = append(items, *CommentItemFromModel(&comments[i], content, favoriteCount))
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
