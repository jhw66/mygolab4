// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"context"
	"errors"
	"log"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageLogic {
	return &GetMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageLogic) GetMessage(req *types.GetMessagesReq) (resp *types.GetMessagesResp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}

	room, _ := l.svcCtx.ChatRepo.FindRoomByRoomId(l.ctx, req.RoomID)
	if room == nil {
		return nil, errors.New("房间不存在")
	}

	if member, err := l.svcCtx.ChatRepo.FindMemberById(l.ctx, req.RoomID, user.ID); member == nil || err != nil {
		return nil, errors.New("您没有权限访问该房间或者访问失败")
	}

	messages, err := l.svcCtx.ChatRepo.ListRoomMessages(l.ctx, req.RoomID, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("获取消息失败")
	}
	total, err := l.svcCtx.ChatRepo.CountRoomMessages(l.ctx, req.RoomID)
	if err != nil {
		log.Println("获取消息数量失败", err)
		total = 0
	}

	return serializer.ChatMessageListRspFromModels(messages, total, req.Page, req.PageSize), nil
}
