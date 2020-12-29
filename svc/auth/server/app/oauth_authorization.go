package app

import (
	"context"
	"time"

	"dfl/lib/cher"
	"dfl/lib/ptr"
	"dfl/svc/auth"
)

func (a *App) Authorization(ctx context.Context, req *auth.AuthorizationRequest, user *auth.User) (*auth.AuthorizationResponse, error) {
	if user.Password == nil || user.InviteCode != nil {
		return nil, cher.New("pending_invite", nil)
	}

	if !checkPasswordHash(req.Password, *user.Password) {
		return nil, cher.New(cher.Unauthorized, nil)
	}

	authCode, expiresAt, err := a.DB.Q.CreateAuthorizationCode(ctx, user.ID, req)
	if err != nil {
		return nil, err
	}

	err = a.DB.Q.ExpireAuthorizationCodes(ctx, ptr.String(authCode))
	if err != nil {
		return nil, err
	}

	expiresIn := expiresAt.Sub(time.Now())

	return &auth.AuthorizationResponse{
		AuthorizationCode: authCode,
		ExpiresAt:         expiresAt.Format(time.RFC3339),
		ExpiresIn:         int(expiresIn.Seconds()),
		State:             req.State,
	}, nil
}
