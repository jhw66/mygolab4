// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package totp

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/twofa"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type TotpEnrollStartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpEnrollStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpEnrollStartLogic {
	return &TotpEnrollStartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpEnrollStartLogic) TotpEnrollStart(req *types.TotpEnrollStartReq) (resp *types.TotpEnrollStartResp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}

	if user.TotpEnabled {
		return nil, errors.New("已启用2FA")
	}

	enrollment, err := twofa.GenerateEnrollment(
		l.svcCtx.Config.Totp.Issuer,
		user.UserName,
		l.svcCtx.Config.Totp.QRCodeLevel,
		l.svcCtx.Config.Totp.QRCodeSize,
	)
	if err != nil {
		return nil, err
	}

	encryptedSecret, err := utils.EncryptString(l.svcCtx.Config.Totp.SecretCipherKey, enrollment.Secret)
	if err != nil {
		return nil, err
	}

	user.TotpPendingSecret = encryptedSecret
	if err := l.svcCtx.UserRepo.Update(l.ctx, user); err != nil {
		return nil, err
	}

	return &types.TotpEnrollStartResp{
		Status: 200,
		Data: &types.TotpEnrollStartItem{
			Issuer:        enrollment.Issuer,
			AccountName:   enrollment.AccountName,
			OTPAuthURL:    enrollment.OtpAuthUrl,
			QRCodeDataURL: enrollment.QrCodeDataUrl,
		},
		Msg: "启动2FA",
	}, nil
}
