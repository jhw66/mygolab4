// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

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

type TotpVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpVerifyLogic {
	return &TotpVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpVerifyLogic) TotpVerify(req *types.TotpVerifyReq, challengeToken string) (resp *types.TotpVerifyResp, err error) {
	challengeToken = strings.TrimSpace(challengeToken)
	code := strings.TrimSpace(req.Code)
	if challengeToken == "" || code == "" {
		return nil, errors.New("challenge token and code are required")
	}

	claim, err := auth.ParseToken(l.svcCtx.Config.Jwt.ChallengeTokenSecret, challengeToken, auth.ChallengeTokenType)
	if err != nil {
		return nil, errors.New("invalid challenge token")
	}

	user, err := l.svcCtx.UserRepo.FindByID(l.ctx, claim.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.TotpEnabled || strings.TrimSpace(user.TotpSecret) == "" {
		return nil, errors.New("2fa is not enabled for this account")
	}

	secret, err := utils.DecryptString(l.svcCtx.Config.Totp.SecretCipherKey, user.TotpSecret)
	if err != nil {
		return nil, err
	}
	if !twofa.ValidateCode(secret, code) {
		return nil, errors.New("invalid totp code")
	}

	accessToken, err := auth.GenerateToken(
		l.svcCtx.Config.Jwt.AccessTokenSecret,
		l.svcCtx.Config.Jwt.AccessTokenExpire,
		user.ID,
		auth.AccessTokenType,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}
	refreshToken, err := auth.GenerateToken(
		l.svcCtx.Config.Jwt.RefreshTokenSecret,
		l.svcCtx.Config.Jwt.RefreshTokenExpire,
		user.ID,
		auth.RefreshTokenType,
	)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &types.TotpVerifyResp{
		Status: 200,
		Data: &types.TotpVerifyItem{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		Msg: "2FA验证成功",
	}, nil
}
