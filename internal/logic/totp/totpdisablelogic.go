// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package totp

import (
	"context"
	"errors"
	"strings"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/twofa"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type TotpDisableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpDisableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpDisableLogic {
	return &TotpDisableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpDisableLogic) TotpDisable(req *types.Disable2FAReq) (resp *types.Disable2FAResp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}

	password := req.Password
	code := strings.TrimSpace(req.Code)
	if password == "" || code == "" {
		return nil, errors.New("密码和验证码不能为空")
	}
	if !utils.ComparePassword(user.PasswordDigest, password) {
		return nil, errors.New("密码错误")
	}
	if !user.TotpEnabled || strings.TrimSpace(user.TotpSecret) == "" {
		return nil, errors.New("2FA未启用")
	}

	secret, err := utils.DecryptString(l.svcCtx.Config.Totp.SecretCipherKey, user.TotpSecret)
	if err != nil {
		return nil, err
	}
	if !twofa.ValidateCode(secret, code) {
		return nil, errors.New("验证码无效")
	}

	user.TotpEnabled = false
	user.TotpSecret = ""
	user.TotpPendingSecret = ""
	if err := l.svcCtx.UserRepo.Update(l.ctx, user); err != nil {
		return nil, err
	}
	return &types.Disable2FAResp{
		Status: 200,
		Data: &types.Disable2FAItem{
			Enabled: false,
		},
		Msg: "2FA已禁用",
	}, nil
}
