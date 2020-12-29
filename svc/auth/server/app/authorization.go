package app

import (
	"context"
	"time"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/ptr"
	"dfl/svc/auth"
)

func (a *App) Authorization(ctx context.Context, req *auth.AuthorizationRequest) (*auth.AuthorizationResponse, error) {
	_, err := a.DB.Q.FindClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	user, err := a.DB.Q.GetUserByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user.Password == nil || user.InviteCode != nil {
		return nil, cher.New("pending_invite", nil)
	}

	if !checkPasswordHash(req.Password, *user.Password) {
		return nil, cher.New(cher.Unauthorized, nil)
	}

	if !authlib.Can(req.Scope, user.Scopes) {
		return nil, cher.New(cher.AccessDenied, nil)
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
