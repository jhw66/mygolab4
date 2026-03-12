package serializer

import "github.com/jhw66/myvideo_lab4/model"

type Comment struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	VideoID   uint   `json:"video_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	User      User   `json:"users"`
}

type CommentList struct {
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Comments []Comment `json:"comments"`
}

func BuildComment(comment *model.Comment) *Comment {
	return &Comment{
		ID:        comment.ID,
		UserID:    comment.UserID,
		VideoID:   comment.VideoID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Unix(),
		User:      *BuildUser(&comment.User),
	}
}

func BuildCommentList(comments *[]model.Comment, total int64, page int, pageSize int) *CommentList {
	var res []Comment
	for key, _ := range *comments {
		res = append(res, *BuildComment(&(*comments)[key]))
	}
	return &CommentList{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Comments: res,
	}
}
