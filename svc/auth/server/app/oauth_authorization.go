package app

import (
	"context"
	"time"

	"dfl/lib/cher"
	"dfl/lib/ptr"
	"dfl/svc/auth"

	"golang.org/x/crypto/bcrypt"
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

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *App) AuthCodeNoNonceExists(ctx context.Context, nonce string) error {
	_, err := a.DB.Q.GetAuthorizationCodeByNonce(ctx, nonce)
	if err != nil {
		if v, ok := err.(cher.E); ok {
			if v.Code == cher.NotFound {
				return nil
			}
		}
	}

	return cher.New("nonce_used", nil)
}
