// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
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
	if err := utils.ValidateRuneLength(req.NickName, 2, 30, "昵称长度需在2-30个字符"); err != nil {
		return &types.UserRsp{Status: 400, Msg: err.Error()}, err
	}
	if err := utils.ValidateRuneLength(req.UserName, 5, 30, "用户名长度需在5-30个字符"); err != nil {
		return &types.UserRsp{Status: 400, Msg: err.Error()}, err
	}
	if err := utils.ValidateRuneLength(req.Password, 8, 40, "密码长度需在8-40个字符"); err != nil {
		return &types.UserRsp{Status: 400, Msg: err.Error()}, err
	}
	if err := utils.ValidateRuneLength(req.PasswordConfirm, 8, 40, "确认密码长度需在8-40个字符"); err != nil {
		return &types.UserRsp{Status: 400, Msg: err.Error()}, err
	}

	if req.PasswordConfirm != req.Password {
		return &types.UserRsp{Status: 400, Msg: "两次输入密码不一致"}, errors.New("两次输入密码不一致")
	}

	exists, Eerr := l.svcCtx.UserRepo.ExistsByUsername(l.ctx, req.UserName)
	if Eerr != nil {
		return &types.UserRsp{Status: 500, Msg: "查询用户失败"}, errors.New("查询用户失败")
	}

	if exists {
		return &types.UserRsp{Status: 409, Msg: "用户名已存在"}, errors.New("用户名已存在")
	}

	exists, Eerr = l.svcCtx.UserRepo.ExistsByNickname(l.ctx, req.NickName)
	if Eerr != nil {
		return &types.UserRsp{Status: 500, Msg: "查询用户失败"}, errors.New("查询用户失败")
	}
	if exists {
		return &types.UserRsp{Status: 409, Msg: "昵称已被占用"}, errors.New("昵称已被占用")
	}

	pwd, hashErr := utils.HashPassword(req.Password)
	if hashErr != nil {
		return &types.UserRsp{Status: 500, Msg: "加密失败"}, errors.New("加密失败")
	}

	user := &model.User{
		NickName:       req.NickName,
		UserName:       req.UserName,
		PasswordDigest: pwd,
	}
	createdUser, cErr := l.svcCtx.UserRepo.Create(l.ctx, user)

	if cErr != nil {
		return &types.UserRsp{Status: 500, Msg: "用户创建失败"}, errors.New("用户创建失败")
	}
	return serializer.UserRspFromModel(createdUser), nil
}
