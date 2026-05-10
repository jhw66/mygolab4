// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *types.UserRsp, err error) {
	if err := utils.ValidateRuneLength(req.NickName, 2, 30); err != nil {
		return nil, err
	}
	if err := utils.ValidateRuneLength(req.UserName, 5, 30); err != nil {
		return nil, err
	}
	if err := utils.ValidateRuneLength(req.Password, 8, 40); err != nil {
		return nil, err
	}
	if err := utils.ValidateRuneLength(req.PasswordConfirm, 8, 40); err != nil {
		return nil, err
	}

	if req.PasswordConfirm != req.Password {
		return nil, errors.New("两次输入密码不一致")
	}

	exists, err := l.svcCtx.UserRepo.ExistsByUsername(l.ctx, req.UserName)
	if err != nil {
		return nil, errors.New("查询用户失败")
	}

	if exists {
		return nil, errors.New("用户名已存在")
	}

	exists, err = l.svcCtx.UserRepo.ExistsByNickname(l.ctx, req.NickName)
	if err != nil {
		return nil, errors.New("查询用户失败")
	}
	if exists {
		return nil, errors.New("昵称已被占用")
	}

	pwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("加密失败")
	}

	user := &model.User{
		NickName:          req.NickName,
		UserName:          req.UserName,
		PasswordDigest:    pwd,
		TotpSecret:        "",
		TotpPendingSecret: "",
		TotpEnabled:       false,
	}
	createdUser, err := l.svcCtx.UserRepo.Create(l.ctx, user)

	if err != nil {
		return nil, errors.New("用户创建失败")
	}
	return &types.UserRsp{
		Status: 200,
		Msg:    "注册成功",
		Data: types.UserItem{
			Id:        createdUser.ID,
			Username:  createdUser.UserName,
			Nickname:  createdUser.NickName,
			CreatedAt: int64(createdUser.CreatedAt.Unix()),
			Avatar:    createdUser.Avatar,
		},
	}, nil
}
