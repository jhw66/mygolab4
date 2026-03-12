package serializer

import "github.com/jhw66/myvideo_lab4/model"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

func BuildUserResponse(user *model.User) *Response {
	return &Response{
		Status: 200,
		Data:   BuildUser(user),
	}
}

func BuildVideoResponse(video *model.Video) *Response {
	return &Response{
		Status: 200,
		Data:   BuildVideo(video),
	}
}

func BuildVideoListResponse(videos *[]model.Video) *Response {
	return &Response{
		Status: 200,
		Data:   BuildVideoList(videos),
	}
}

func BuildCommentResponse(comment *model.Comment) *Response {
	return &Response{
		Status: 200,
		Data:   BuildComment(comment),
	}
}

func BuildCommentListResponse(comments *[]model.Comment, total int64, page int, pageSize int) *Response {
	return &Response{
		Status: 200,
		Data:   BuildCommentList(comments, total, page, pageSize),
	}
}
