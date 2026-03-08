package serializer

import "github.com/jhw66/myvideo_lab4/model"

type User struct {
	ID        uint   `json:"id"`
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	CreatedAt int64  `json:"created_at"`
	Avatar    string `json:"avatar"`
}

func BuildUser(user *model.User) *User {
	return &User{
		ID:        user.ID,
		UserName:  user.UserName,
		NickName:  user.NickName,
		CreatedAt: user.CreatedAt.Unix(),
		Avatar:    user.Avatar,
	}
}
