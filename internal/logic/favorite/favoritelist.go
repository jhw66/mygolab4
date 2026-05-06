// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package favorite

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList() (resp *types.VideoListRsp, err error) {
	user, ok := utlcontext.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.VideoListRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	videos, err := l.svcCtx.FavoriteRepo.ListVideosByUserID(l.ctx, user.ID)
	if err != nil {
		return &types.VideoListRsp{Status: 500, Msg: "查询失败"}, errors.New("查询失败")
	}
	return serializer.VideoListRspFromModels(videos), nil
}
