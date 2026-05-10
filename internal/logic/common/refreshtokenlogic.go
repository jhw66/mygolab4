// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	claims, err := auth.ParseToken(l.svcCtx.Config.Jwt.RefreshTokenSecret, req.RefreshToken, auth.RefreshTokenType)
	if err != nil {
		return nil, errors.New("无效的refresh令牌")
	}

	token, err := auth.GenerateToken(l.svcCtx.Config.Jwt.AccessTokenSecret,
		l.svcCtx.Config.Jwt.AccessTokenExpire,
		claims.UserID,
		auth.AccessTokenType)
	if err != nil {
		return nil, errors.New("生成新accesstoken失败")
	}

	return &types.RefreshTokenResp{
		Status: 200,
		Msg:    "刷新成功",
		Data: types.RefreshTokenItem{
			AccessToken: token,
		},
	}, nil
}
