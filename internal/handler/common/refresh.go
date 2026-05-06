// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"
	"time"

	"github.com/jhw66/myvideo_lab4/internal/logic/common"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshCookie, cookieErr := r.Cookie("refresh_token")
		if cookieErr != nil || refreshCookie.Value == "" {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 401,
				Msg:    "未找到刷新令牌，请重新登录",
			})
			return
		}

		l := common.NewRefreshLogic(r.Context(), svcCtx)
		resp, newAccessToken, err := l.Refresh(refreshCookie.Value)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
		} else {
			if newAccessToken != "" {
				http.SetCookie(w, &http.Cookie{
					Name:   "access_token",
					Value:  newAccessToken,
					Path:   "/",
					MaxAge: int((15 * time.Minute).Seconds()),
				})
			}
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
		}
	}
}
