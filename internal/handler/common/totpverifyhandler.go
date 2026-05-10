// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/common"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TotpVerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TotpVerifyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, types.CommonRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}

		challengeToken, err := auth.GetTokenFromHeader(r)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, types.CommonRsp{
				Status: 400,
				Msg:    "未传challengeToken",
				Error:  err.Error(),
			})
			return
		}

		l := common.NewTotpVerifyLogic(r.Context(), svcCtx)
		resp, err := l.TotpVerify(&req, challengeToken)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, types.CommonRsp{
				Status: 400,
				Msg:    "验证2FA失败",
				Error:  err.Error(),
			})
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    resp.Data.RefreshToken,
				Path:     "/api/v1/refresh",
				MaxAge:   int(svcCtx.Config.Jwt.RefreshTokenExpire),
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    resp.Data.AccessToken,
				Path:     "/",
				MaxAge:   int(svcCtx.Config.Jwt.AccessTokenExpire),
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
			})
			resp.Data.AccessToken = ""
			resp.Data.RefreshToken = ""
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
