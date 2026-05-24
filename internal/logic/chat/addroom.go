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
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/chat"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoomLogic {
	return &AddRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRoomLogic) AddRoom(req *types.AddRoomReq) (resp *types.AddRoomResp, err error) {
	if err := utils.ValidateRuneLength(req.RoomName, 5, 20); err != nil {
		return nil, err
	}
	req.Visibility = strings.ToLower(req.Visibility)
	if err := utils.ValidateOptional(req.Visibility, chat.VisibilityPublic, chat.VisibilityPrivate); err != nil {
		return nil, err
	}
	if req.Visibility == chat.VisibilityPrivate {
		if err := utils.ValidateRuneLength(req.Password, 8, 40); err != nil {
			return nil, err
		}
	}

	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}

	if room, _ := l.svcCtx.ChatRepo.FindRoomByRoomName(l.ctx, req.RoomName); room != nil {
		return nil, errors.New("房间名称已存在")
	}

	room := &model.ChatRoom{
		RoomName:   req.RoomName,
		Visibility: req.Visibility,
		OwnerID:    user.ID,
	}
	if req.Visibility == chat.VisibilityPrivate {
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("加密失败")
		}
		room.PasswordHash = passwordHash
	}

	err = l.svcCtx.TransactionRepository.WithTransaction(l.ctx, func(tx *gorm.DB) error {
		if err := l.svcCtx.ChatRepo.CreateRoomWithTx(l.ctx, tx, room); err != nil {
			return err
		}
		return l.svcCtx.ChatRepo.JoinRoomWithTx(l.ctx, tx, room.ID, user.ID, chat.RoleOwner)
	})
	if err != nil {
		return nil, errors.New("创建房间失败")
	}

	return &types.AddRoomResp{
		Status: 200,
		Data:   serializer.ChatRoomItemFromModel(room),
		Msg:    "创建房间成功",
		Error:  "",
	}, nil
}
