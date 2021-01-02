package app

import (
	"dfl/lib/cache"
	"dfl/svc/short/server/db"
	"dfl/svc/short/server/lib/storageproviders"

	"github.com/sirupsen/logrus"
	hashids "github.com/speps/go-hashids"
)

// App is a struct for the App and it's handlers
type App struct {
	Logger *logrus.Logger

	DB     *db.DB
	SP     storageproviders.StorageProvider
	Hasher *hashids.HashID
	Redis  *cache.Cache

	RootURL string
}
