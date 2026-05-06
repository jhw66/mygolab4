package middleware

// const (
// 	defaultSessionName   = "my_cookie"
// 	defaultSessionSecret = "secret"
// )

// type SessionAccountMiddleware struct {
// 	store    sessions.Store
// 	UserRepo userrepo.UserRepository
// }

// func NewSessionAccountMiddleware(secret string, userRepo userrepo.UserRepository) *SessionAccountMiddleware {
// 	if secret == "" {
// 		secret = defaultSessionSecret
// 	}
// 	return &SessionAccountMiddleware{
// 		store:    sessions.NewCookieStore([]byte(secret)),
// 		UserRepo: userRepo,
// 	}
// }

// // Handle 会尝试从 session 中读取 user_id，并在读取成功时注入 context 用户信息。
// // 若 session 不存在或无效，不会中断请求，便于与 AuthLoginMiddleware 组合使用。
// func (m *SessionAccountMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, err := m.store.Get(r, defaultSessionName)
// 		if err == nil {
// 			if rawID := session.Values["user_id"]; rawID != nil {
// 				if userID, ok := rawID.(string); ok && userID != "" {
// 					if user, err := m.UserRepo.FindByID(r.Context(), userID); err == nil && user != nil {
// 						ctx := context.WithValue(r.Context(), UserContextKey, user)
// 						r = r.WithContext(ctx)
// 					}
// 				}
// 			}
// 		}
// 		next(w, r)
// 	}
// }

// type AuthLoginMiddleware struct {
// }

// func NewAuthLoginMiddleware() *AuthLoginMiddleware {
// 	return &AuthLoginMiddleware{}
// }

// // Handle 校验 context 中是否有已登录用户。
// func (m *AuthLoginMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		user, ok := getUserFromContext(r.Context())
// 		if !ok || user == nil {
// 			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, serializer.Response{
// 				Status: http.StatusUnauthorized,
// 				Msg:    "需要登录",
// 			})
// 			return
// 		}

// 		next(w, r)
// 	}
// }
// func getUserFromContext(ctx context.Context) (*model.User, bool) {
// 	user, ok := ctx.Value(UserContextKey).(*model.User)
// 	if !ok {
// 		return nil, false
// 	}
// 	return user, true
// }
