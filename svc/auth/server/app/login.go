package app

import (
	"context"

	"dfl/lib/cher"
	"dfl/svc/auth"
)

func (a *App) Login(ctx context.Context, user *auth.User, password string) error {
	if user.Password == nil || user.InviteCode != nil {
		return cher.New("pending_invite", nil)
	}

	if !checkPasswordHash(password, *user.Password) {
		return cher.New(cher.Unauthorized, nil)
	}

	return nil
}
