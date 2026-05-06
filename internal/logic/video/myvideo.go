// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMyVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyVideoLogic {
	return &MyVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyVideoLogic) MyVideo() (resp *types.VideoListRsp, err error) {
	user, ok := utlcontext.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.VideoListRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}

	videos, err := l.svcCtx.VideoRepo.FindByUserID(l.ctx, user.ID)
	if err != nil {
		return &types.VideoListRsp{Status: 500, Msg: "查找失败"}, errors.New("查找失败")
	}

	return serializer.VideoListRspFromModels(videos), nil
}
