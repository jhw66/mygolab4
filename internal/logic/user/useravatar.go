// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAvatarLogic {
	return &UserAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAvatarLogic) UserAvatar(r *http.Request) (resp *types.UserRsp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}
	if err = utils.ParseMultipartForm(r, 32<<20); err != nil {
		return nil, errors.New("请求格式错误")
	}
	file, err := utils.GetFormFile(r, "avatar")
	if err != nil {
		return nil, errors.New("头像文件不能为空")
	}

	avatarDiskPath, avatarWebPath, cleanupOldAvatar, err := utils.ReplaceStoredFile("static/avatar", "avatar", user.ID, file, user.Avatar)
	if err != nil {
		return nil, errors.New("保存头像失败")
	}

	user.Avatar = avatarWebPath
	err = l.svcCtx.UserRepo.Update(l.ctx, user)
	if err != nil {
		if err := utils.RemoveIfExists(avatarDiskPath); err != nil {
			l.Errorf("failed to remove new avatar file after update failure, path=%s, err=%v", avatarDiskPath, err)
		}
		return nil, errors.New("更新用户头像失败")
	}
	if cleanupOldAvatar != nil {
		if err := cleanupOldAvatar(); err != nil {
			l.Errorf("cleanup old avatar failed, err=%v", err)
		}
	}

	return serializer.UserRspFromModel(user), nil
}
