package app

import (
	"context"
	"time"

	"dfl/lib/cher"
	"dfl/svc/auth"

	"golang.org/x/crypto/bcrypt"
)

func (a *App) Register(ctx context.Context, req *auth.RegisterRequest) error {
	user, err := a.DB.Q.GetUserByName(ctx, req.Username)
	if err != nil {
		return err
	}

	if user.InviteCode == nil || user.InviteExpiry == nil || user.Password != nil {
		return cher.New("invite_accepted", nil)
	}

	if user.InviteCode != nil && *user.InviteCode != req.InviteCode {
		return cher.New("bad_code", nil)
	}

	if user.InviteExpiry.Before(time.Now()) {
		return cher.New("invite_expired", nil)
	}

	hashedPW, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	err = a.DB.Q.RedeemInvite(ctx, user.ID, hashedPW)
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
