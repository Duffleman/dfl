package app

import (
	"dfl/lib/templates"
	"dfl/svc/auth/server/db"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger *logrus.Entry

	WA *webauthn.WebAuthn
	DB *db.DB
	SK *SigningKeys

	Template *templates.Template

	JWTIssuer string
}
