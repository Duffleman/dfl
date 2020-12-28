package config

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger *logrus.Logger

	DB *sql.DB
}
