// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshLogic) Refresh(refreshToken string) (resp *types.CommonRsp, newAccessToken string, err error) {
	claims, parseErr := utils.ParseToken(refreshToken)
	if parseErr != nil {
		return &types.CommonRsp{
			Status: 401,
			Msg:    parseErr.Error(),
		}, "", errors.New(parseErr.Error())
	}
	if claims.TokenType != "refresh" {
		return &types.CommonRsp{
			Status: 401,
			Msg:    "refresh令牌类型错误",
		}, "", errors.New("refresh令牌类型错误")
	}

	token, genErr := utils.GenerateAccessToken(claims.UserID)
	if genErr != nil {
		return &types.CommonRsp{
			Status: 500,
			Msg:    "生成新accesstoken失败",
		}, "", errors.New("生成新accesstoken失败")
	}

	return &types.CommonRsp{
		Status: 200,
		Msg:    "重新生成accesstoken成功",
	}, token, nil
}
