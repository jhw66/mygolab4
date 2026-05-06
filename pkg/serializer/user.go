package serializer

import (
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
)

type User struct {
	ID        string `json:"id"`
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	CreatedAt int64  `json:"created_at"`
	Avatar    string `json:"avatar"`
}

func UserItemFromModel(user *model.User) types.UserItem {
	return types.UserItem{
		Id:        user.ID,
		Username:  user.UserName,
		Nickname:  user.NickName,
		CreatedAt: user.CreatedAt.Unix(),
		Avatar:    user.Avatar,
	}
}

func UserRspFromModel(user *model.User) *types.UserRsp {
	return &types.UserRsp{
		Status: 200,
		Data:   UserItemFromModel(user),
	}
}

func UserListRspFromModels(users []model.User) *types.UserListRsp {
	items := make([]types.UserItem, 0, len(users))
	for i := range users {
		items = append(items, UserItemFromModel(&users[i]))
	}
	return &types.UserListRsp{
		Status: 200,
		Data:   items,
	}
}
