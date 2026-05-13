// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package relation

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationReq) (resp *types.CommonRsp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}
	if user.ID == req.Uid {
		return nil, errors.New("不可关注自己")
	}
	if _, err := l.svcCtx.UserRepo.FindByID(l.ctx, req.Uid); err != nil {
		return nil, errors.New("目标用户不存在")
	}

	exists, err := l.svcCtx.RelationRepo.ExistsByUserAndTarget(l.ctx, user.ID, req.Uid)
	if err != nil {
		return nil, errors.New("关注操作失败")
	}
	if exists {
		if err := l.svcCtx.RelationRepo.DeleteByUserAndTarget(l.ctx, user.ID, req.Uid); err != nil {
			return nil, errors.New("取消关注失败")
		}
		return nil, nil
	}
	if err := l.svcCtx.RelationRepo.Create(l.ctx, user.ID, req.Uid); err != nil {
		return nil, errors.New("关注失败")
	}
	return &types.CommonRsp{Status: 200, Msg: "关注成功"}, nil
}
