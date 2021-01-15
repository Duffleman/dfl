package app

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

func (a *App) DeleteKey(ctx context.Context, userID, keyID string) error {
	credential, err := a.DB.Q.FindU2FCredential(ctx, keyID)
	if err != nil {
		return err
	}

	if credential.UserID != userID {
		return cher.New(cher.AccessDenied, nil)
	}

	allCredentials, err := a.ListU2FKeys(ctx, userID, false)
	if err != nil {
		return err
	}

	if len(allCredentials) == 1 && keyID == allCredentials[0].ID {
		return cher.New("last_signed_key", cher.M{"key_id": keyID})
	}

	return a.DB.Q.DeleteU2FCredential(ctx, userID, keyID)
}
