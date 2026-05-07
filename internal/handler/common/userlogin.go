// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"
	"time"

	"github.com/jhw66/myvideo_lab4/internal/logic/common"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/utils"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.UserRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}

		l := common.NewUserLoginLogic(r.Context(), svcCtx)
		resp, err := l.UserLogin(&req)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
		} else {
			if resp.Data.Id != "" {
				accessToken, accessErr := utils.GenerateAccessToken(resp.Data.Id)
				refreshToken, refreshErr := utils.GenerateRefreshToken(resp.Data.Id)
				if accessErr == nil && refreshErr == nil {
					http.SetCookie(w, &http.Cookie{
						Name:     "refresh_token",
						Value:    refreshToken,
						Path:     "/api/v1/refresh",
						MaxAge:   int((24 * time.Hour).Seconds()),
						SameSite: http.SameSiteStrictMode,
						HttpOnly: true,
					})
					http.SetCookie(w, &http.Cookie{
						Name:   "access_token",
						Value:  accessToken,
						Path:   "/",
						MaxAge: int((15 * time.Minute).Seconds()),
					})
				} else {
					httpx.WriteJsonCtx(r.Context(), w, 200, types.CommonRsp{
						Status: 500,
						Msg:    "登录成功，但生成令牌失败",
					})
					return
				}
			}
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
		}
	}

}
