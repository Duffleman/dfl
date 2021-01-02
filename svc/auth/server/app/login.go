package app

import (
	"context"

	"dfl/lib/cher"
	"dfl/svc/auth"
)

func (a *App) CheckLoginValidity(ctx context.Context, user *auth.User) error {
	if user.InviteRedeemedAt == nil {
		return cher.New("pending_invite", nil)
	}

	return nil
}
