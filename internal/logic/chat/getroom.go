// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoomLogic {
	return &GetRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoomLogic) GetRoom() (resp *types.GetRoomResp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}

	room, err := l.svcCtx.ChatRepo.ListVisibleRooms(l.ctx, user.ID)
	if err != nil {
		return nil, errors.New("获取房间列表失败")
	}

	return serializer.ChatRoomListRspFromModels(room), nil
}
