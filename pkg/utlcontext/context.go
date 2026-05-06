package utlcontext

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUserInContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}
