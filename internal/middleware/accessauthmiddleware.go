// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/pkg/db/repository/user"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AccessAuthMiddleware struct {
	UserRepo user.UserRepository
}

func NewAccessAuthMiddleware(userRepo user.UserRepository) *AccessAuthMiddleware {
	return &AccessAuthMiddleware{
		UserRepo: userRepo,
	}
}

func (m *AccessAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, serializer.Response{
				Status: http.StatusUnauthorized,
				Msg:    "无access令牌",
			})
			return
		}

		claims, err := utils.ParseToken(cookie.Value)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, serializer.Response{
				Status: http.StatusUnauthorized,
				Msg:    "access令牌过期或者不合法",
			})
			return
		}

		if claims.TokenType != "access" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, serializer.Response{
				Status: http.StatusUnauthorized,
				Msg:    "access令牌类型错误",
			})
			return
		}

		user, err := m.UserRepo.FindByID(r.Context(), claims.UserID)
		if err != nil || user == nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, serializer.Response{
				Status: http.StatusUnauthorized,
				Msg:    "用户不存在或者access令牌无效",
			})
			return
		}

		next(w, r.WithContext(utlcontext.WithUserInContext(r.Context(), user)))
	}
}
