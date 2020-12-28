package app

import (
	"dfl/svc/auth/server/db"
)

type App struct {
	DB *db.DB
	SK *SigningKeys
}
