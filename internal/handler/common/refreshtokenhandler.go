// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/common"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenReq
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 401,
				Msg:    "未找到刷新令牌，请重新登录",
				Error:  err.Error(),
			})
			return
		}
		req.RefreshToken = refreshCookie.Value

		l := common.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 401,
				Msg:    "刷新失败",
				Error:  err.Error(),
			})
			return
		}
		accessToken := resp.Data.AccessToken
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			MaxAge:   int(svcCtx.Config.Jwt.AccessTokenExpire),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})
		resp.Data = nil
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
