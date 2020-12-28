package auth

import (
	"context"
)

type contextKey string

var UserContextKey contextKey = "user"

type AuthUser struct {
	UserID   string
	Username string
	Scopes   string
}

func (au AuthUser) Can(action string) bool {
	return Can(action, au.Scopes)
}

func GetFromContext(ctx context.Context) AuthUser {
	return ctx.Value(UserContextKey).(AuthUser)
}
