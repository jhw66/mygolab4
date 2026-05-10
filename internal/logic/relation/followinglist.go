// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package relation

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowingListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowingListLogic {
	return &FollowingListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowingListLogic) FollowingList() (resp *types.UserListRsp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.UserListRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	users, err := l.svcCtx.RelationRepo.ListFollowingUsers(l.ctx, user.ID)
	if err != nil {
		return &types.UserListRsp{Status: 500, Msg: "查询关注列表失败"}, errors.New("查询关注列表失败")
	}
	rsp := serializer.UserListRspFromModels(users)
	rsp.Msg = "查询关注列表成功"
	return rsp, nil
}
