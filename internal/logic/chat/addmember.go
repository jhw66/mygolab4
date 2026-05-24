// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"context"
	"errors"
	"strings"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/chat"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMemberLogic {
	return &AddMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMemberLogic) AddMember(req *types.AddMemberReq) (resp *types.AddMemberResp, err error) {
	req.Role = strings.ToLower(req.Role)
	if err := utils.ValidateOptional(req.Role, chat.RoleMember, chat.RoleAdmin); err != nil {
		return nil, err
	}

	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	if room, _ := l.svcCtx.ChatRepo.FindRoomByRoomId(l.ctx, req.RoomID); room == nil {
		return nil, errors.New("房间不存在")
	}

	member, _ := l.svcCtx.ChatRepo.FindMemberById(l.ctx, req.RoomID, user.ID)
	if member == nil {
		return nil, errors.New("您没有权限")
	}
	if err := utils.ValidateOptional(member.Role, chat.RoleOwner, chat.RoleAdmin); err != nil {
		return nil, errors.New("您没有权限")
	}
	if targetUser, _ := l.svcCtx.UserRepo.FindByID(l.ctx, req.UserID); targetUser == nil {
		return nil, errors.New("用户不存在")
	}

	exists, _ := l.svcCtx.ChatRepo.FindMemberById(l.ctx, req.RoomID, req.UserID)
	if exists != nil {
		return nil, errors.New("该用户已加入房间")
	}

	if err := l.svcCtx.ChatRepo.JoinRoom(l.ctx, req.RoomID, req.UserID, req.Role); err != nil {
		return nil, errors.New("加入房间失败")
	}
	return &types.AddMemberResp{
		Status: 200,
		Msg:    "加入房间成功",
		Error:  "",
	}, nil
}
