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

type TotpEnrollConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpEnrollConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpEnrollConfirmLogic {
	return &TotpEnrollConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpEnrollConfirmLogic) TotpEnrollConfirm(req *types.TotpEnrollConfirmReq) (resp *types.TotpEnrollConfirmResp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, errors.New("验证码不能为空")
	}

	secret, err := utils.DecryptString(l.svcCtx.Config.Totp.SecretCipherKey, user.TotpPendingSecret)
	if err != nil {
		return nil, err
	}
	if !twofa.ValidateCode(secret, code) {
		return nil, errors.New("验证码无效")
	}

	user.TotpSecret = user.TotpPendingSecret
	user.TotpPendingSecret = ""
	user.TotpEnabled = true
	if err := l.svcCtx.UserRepo.Update(l.ctx, user); err != nil {
		return nil, err
	}

	return &types.TotpEnrollConfirmResp{
		Status: 200,
		Data: &types.TotpEnrollConfirmItem{
			Enabled: true,
		},
		Msg: "确认2FA启用成功",
	}, nil

}
