package app

import (
	"context"
	"time"

	"dfl/lib/cher"
	"dfl/svc/auth"
)

func (a *App) Register(ctx context.Context, user *auth.User) error {
	if err := a.CheckRegistrationValidity(ctx, user, nil); err != nil {
		return err
	}

	return a.DB.Q.RedeemInvite(ctx, user.ID)
}

func (a *App) CheckRegistrationValidity(ctx context.Context, user *auth.User, inviteCode *string) error {
	if user.InviteRedeemedAt != nil {
		return cher.New("already_registered", nil)
	}

	if user.InviteExpiry.Before(time.Now()) {
		return cher.New("invite_expired", nil)
	}

	if inviteCode != nil {
		if user.InviteCode != *inviteCode {
			return cher.New("bad_invite_code", nil)
		}
	}

	return nil
}
