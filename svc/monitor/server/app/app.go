package app

import (
	"net/http"

	"dfl/svc/monitor"
	"dfl/svc/monitor/server/lib/cachet"

	"github.com/sirupsen/logrus"
)

// App is a struct for the app methods to attach to
type App struct {
	Logger *logrus.Logger

	Client           *http.Client
	ClientNoValidate *http.Client
	Jobs             []*monitor.Job
	Cachet           *cachet.Client
}
