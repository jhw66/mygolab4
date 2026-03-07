package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/utils"
)

type UserLogin struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
}

func (userlogin UserLogin) Login() (model.User, *serializer.Response) {
	var user model.User
	if err := model.Db.Where("user_name = ?", userlogin.UserName).Take(&user).Error; err != nil {
		return user, &serializer.Response{
			Status: 404,
			Msg:    "用户不存在，请先注册",
		}
	}

	if !utils.ComparePassword(user.PasswordDigest, userlogin.Password) {
		return user, &serializer.Response{
			Status: 403,
			Msg:    "密码错误",
		}
	}

	return user, nil
}
