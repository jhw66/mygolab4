// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/user"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserLogoutLogic(r.Context(), svcCtx)
		resp, err := l.UserLogout()
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/api/v1/refresh",
				MaxAge:   -1,
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
			})
			http.SetCookie(w, &http.Cookie{
				Name:   "access_token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
		}
	}
}
