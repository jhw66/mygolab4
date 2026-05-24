// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/chat"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinRoomLogic {
	return &JoinRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinRoomLogic) JoinRoom(req *types.JoinRoomReq) (resp *types.JoinRoomResp, err error) {
	if req.Password != "" {
		if err := utils.ValidateRuneLength(req.Password, 8, 40); err != nil {
			return nil, err
		}
	}

	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}
	room, _ := l.svcCtx.ChatRepo.FindRoomByRoomId(l.ctx, req.RoomID)
	if room == nil {
		return nil, errors.New("房间不存在")
	}
	exists, _ := l.svcCtx.ChatRepo.FindMemberById(l.ctx, req.RoomID, user.ID)
	if exists != nil {
		return nil, errors.New("您已加入房间")
	}

	if room.Visibility == chat.VisibilityPrivate {
		if !utils.ComparePassword(room.PasswordHash, req.Password) {
			return nil, errors.New("密码错误")
		}
		if err := l.svcCtx.ChatRepo.JoinRoom(l.ctx, req.RoomID, user.ID, chat.RoleMember); err != nil {
			return nil, errors.New("加入房间失败")
		}

	} else {
		if err := l.svcCtx.ChatRepo.JoinRoom(l.ctx, req.RoomID, user.ID, chat.RoleMember); err != nil {
			return nil, errors.New("加入房间失败")
		}
	}

	return &types.JoinRoomResp{
		Status: 200,
		Msg:    "加入房间成功",
		Error:  "",
	}, nil
}
