package auth

import (
	"context"
	m "foodapp/internal/models"
)

type contextKey struct{}

var UserContextKey contextKey

func AddUserContext(ctx context.Context, user m.UserContext) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

func GetUserContext(ctx context.Context) (m.UserContext, bool) {
	userCtx, ok := ctx.Value(UserContextKey).(m.UserContext)
	return userCtx, ok
}
