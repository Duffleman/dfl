package app

import (
	"dfl/svc/auth/server/db"

	"github.com/duo-labs/webauthn/webauthn"
)

type App struct {
	WA *webauthn.WebAuthn
	DB *db.DB
	SK *SigningKeys
}
