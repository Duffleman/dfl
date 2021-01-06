package app

import (
	"net"
	"time"

	"dfl/svc/monitor"

	sdk "github.com/andygrunwald/cachet"
	log "github.com/sirupsen/logrus"
)

// TCPTimeout is the time the TCP has to open before we consider it to have failed
const TCPTimeout = 10 * time.Second

func (a *App) doTCP(job *monitor.Job) int {
	conn, err := net.DialTimeout("tcp", job.Host, TCPTimeout)
	if err != nil {
		log.Infof("cannot open tcp to host %s", job.Host)
		return sdk.ComponentStatusMajorOutage
	}
	defer conn.Close()

	return sdk.ComponentStatusOperational
}
