package serializer

import "github.com/jhw66/myvideo_lab4/model"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

func BuildUserResponse(user *model.User) Response {
	return Response{
		Status: 200,
		Data:   BuildUser(user),
	}
}

func BuildVedioResponse(vedio *model.Vedio) Response {
	return Response{
		Status: 200,
		Data:   BuildVedio(vedio),
	}
}
