// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
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

func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	if err := utils.ValidateRuneLength(req.UserName, 5, 30); err != nil {
		return nil, err
	}
	if err := utils.ValidateRuneLength(req.Password, 8, 40); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserRepo.FindByUsername(l.ctx, req.UserName)
	if err != nil {
		return nil, errors.New("用户不存在，请先注册")
	}

	if !utils.ComparePassword(user.PasswordDigest, req.Password) {
		return nil, errors.New("密码错误")
	}

	if user.TotpEnabled {
		challengeToken, err := auth.GenerateToken(
			l.svcCtx.Config.Jwt.ChallengeTokenSecret,
			l.svcCtx.Config.Jwt.ChallengeTokenExpire,
			user.ID,
			auth.ChallengeTokenType)
		if err != nil {
			return nil, err
		}

		return &types.UserLoginResp{
			Status: 200,
			Data: types.LoginItem{
				Need2FA:        true,
				ChallengeToken: challengeToken,
			},
			Msg: "登录成功，等待二次验证",
		}, nil
	}

	accessToken, err := auth.GenerateToken(l.svcCtx.Config.Jwt.AccessTokenSecret,
		l.svcCtx.Config.Jwt.AccessTokenExpire,
		user.ID,
		auth.AccessTokenType)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateToken(l.svcCtx.Config.Jwt.RefreshTokenSecret,
		l.svcCtx.Config.Jwt.RefreshTokenExpire,
		user.ID,
		auth.RefreshTokenType)
	if err != nil {
		return nil, err
	}

	return &types.UserLoginResp{
		Status: 200,
		Data: types.LoginItem{
			Need2FA:      false,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		Msg: "登录成功",
	}, nil
}
