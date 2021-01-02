package app

import (
	"context"
	"time"

	"dfl/svc/auth"

	"github.com/duo-labs/webauthn/webauthn"
)

func (a *App) CreateU2FCredential(ctx context.Context, user *auth.User, challengeID string, keyName *string, credential *webauthn.Credential, signedAt *time.Time) error {
	return a.DB.Q.CreateU2FCredential(ctx, user, challengeID, keyName, credential, signedAt)
}
