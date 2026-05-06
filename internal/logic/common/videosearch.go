// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoSearchLogic {
	return &VideoSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoSearchLogic) VideoSearch(req *types.VideoSearchReq) (resp *types.VideoListRsp, err error) {
	videos, err := l.svcCtx.VideoRepo.SearchByTitle(l.ctx, req.KeyWord)
	if err != nil {
		return &types.VideoListRsp{Status: 500, Msg: "查找失败"}, errors.New("查找失败")
	}
	if len(videos) == 0 {
		return &types.VideoListRsp{Status: 404, Msg: "未找到相关视频"}, errors.New("未找到相关视频")
	}

	return serializer.VideoListRspFromModels(videos), nil
}
