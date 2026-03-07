package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/utils"
)

type UserRegister struct {
	NickName        string `form:"nick_name" json:"nick_name" binding:"required,min=2,max=30"`
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

func (useregister UserRegister) Register() (model.User, *serializer.Response) {
	var user model.User
	if useregister.PasswordConfirm != useregister.Password {
		return user, &serializer.Response{
			Status: 400,
			Msg:    "两次输入密码不一致",
		}
	}

	var count int64
	if model.Db.Where("user_name = ?", useregister.UserName).Take(&user).Count(&count); count != 0 {
		return user, &serializer.Response{
			Status: 409,
			Msg:    "用户名已存在",
		}
	}

	if model.Db.Where("nick_name = ?", useregister.NickName).Take(&user).Count(&count); count != 0 {
		return user, &serializer.Response{
			Status: 409,
			Msg:    "昵称已被占用",
		}
	}
	user.NickName = useregister.NickName
	user.UserName = useregister.UserName
	passworddigest, err := utils.HashPassword(useregister.Password)
	if err != nil {
		return user, &serializer.Response{
			Status: 500,
			Msg:    "加密失败",
		}
	}
	user.PasswordDigest = passworddigest
	err = model.Db.Create(&user).Error
	if err != nil {
		return user, &serializer.Response{
			Status: 500,
			Msg:    "用户创建失败",
		}
	}
	return user, nil

}
