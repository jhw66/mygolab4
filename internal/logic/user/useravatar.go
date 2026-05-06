// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"

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
	user, ok := utlcontext.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.UserRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	if err = utils.ParseMultipartForm(r, 32<<20); err != nil {
		return &types.UserRsp{Status: 400, Msg: "请求格式错误"}, errors.New("请求格式错误")
	}
	file, err := utils.GetFormFile(r, "avatar")
	if err != nil {
		return &types.UserRsp{Status: 400, Msg: "头像文件不能为空"}, errors.New("头像文件不能为空")
	}

	avatarDiskPath, avatarWebPath, cleanupOldAvatar, err := utils.ReplaceStoredFile("static/avatar", "avatar", user.ID, file, user.Avatar)
	if err != nil {
		return &types.UserRsp{Status: 500, Msg: "保存头像失败"}, errors.New("保存头像失败")
	}

	user.Avatar = avatarWebPath
	err = l.svcCtx.UserRepo.Update(l.ctx, user)
	if err != nil {
		_ = utils.RemoveIfExists(avatarDiskPath)
		return &types.UserRsp{Status: 500, Msg: "更新用户头像失败"}, errors.New("更新用户头像失败")
	}
	if cleanupOldAvatar != nil {
		_ = cleanupOldAvatar()
	}

	return serializer.UserRspFromModel(user), nil
}
