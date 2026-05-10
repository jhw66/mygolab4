// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/config"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/user"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AccessAuthMiddleware struct {
	UserRepo user.UserRepository
	config   config.Config
}

func NewAccessAuthMiddleware(userRepo user.UserRepository, c config.Config) *AccessAuthMiddleware {
	return &AccessAuthMiddleware{
		UserRepo: userRepo,
		config:   c,
	}
}

func (m *AccessAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, types.CommonRsp{
				Status: http.StatusUnauthorized,
				Msg:    "用户未登录",
				Error:  "missing access_token cookie",
			})
			return
		}

		claims, err := auth.ParseToken(
			m.config.Jwt.AccessTokenSecret,
			cookie.Value,
			auth.AccessTokenType,
		)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, types.CommonRsp{
				Status: http.StatusUnauthorized,
				Msg:    "access令牌过期或者不合法",
				Error:  err.Error(),
			})
			return
		}

		user, err := m.UserRepo.FindByID(r.Context(), claims.UserID)
		if err != nil || user == nil {
			errMsg := "user not found"
			if err != nil {
				errMsg = err.Error()
			}
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, types.CommonRsp{
				Status: http.StatusUnauthorized,
				Msg:    "用户不存在或者access令牌无效",
				Error:  errMsg,
			})
			return
		}

		next(w, r.WithContext(auth.WithUserInContext(r.Context(), user)))
	}
}
