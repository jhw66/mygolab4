// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq) (resp *types.UserRsp, err error) {
	user, err := l.svcCtx.UserRepo.FindByUsername(l.ctx, req.UserName)
	if err != nil {
		return &types.UserRsp{Status: 404, Msg: "用户不存在，请先注册"}, errors.New("用户不存在，请先注册")
	}

	if !utils.ComparePassword(user.PasswordDigest, req.Password) {
		return &types.UserRsp{Status: 403, Msg: "密码错误"}, errors.New("密码错误")
	}

	return serializer.UserRspFromModel(user), nil
}
