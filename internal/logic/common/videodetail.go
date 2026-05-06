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
	"gorm.io/gorm"
)

type VideoDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoDetailLogic {
	return &VideoDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoDetailLogic) VideoDetail(req *types.VideoIdReq) (resp *types.VideoRsp, err error) {
	if req.Id == "" {
		return &types.VideoRsp{Status: 400, Msg: "请传入视频id"}, errors.New("请传入视频id")
	}

	video, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.VideoRsp{Status: 404, Msg: "未找到该视频"}, errors.New("未找到该视频")
		}
		return &types.VideoRsp{Status: 500, Msg: "查找失败"}, errors.New("查找失败")
	}

	return serializer.VideoRspFromModel(video), nil
}
