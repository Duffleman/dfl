package app

import (
	"fmt"

	sdk "github.com/andygrunwald/cachet"
	log "github.com/sirupsen/logrus"
)

func (a *App) handleMessage(jw jobWrap) {
	if jw.Job.ComponentID == 0 {
		log.Warnf("cannot update component %s, no matched catchet component", jw.Job.Name)
		return
	}

	_, _, err := a.Cachet.Components.Update(jw.Job.ComponentID, &sdk.Component{
		Status: jw.Outcome,
	})
	if err != nil {
		log.Warn(fmt.Errorf("cannot update component %s: %w", jw.Job.ComponentName, err))
		return
	}
}

func (a *App) messageHandlerWorker(ch chan jobWrap) {
	for {
		select {
		case jw, ok := <-ch:
			if !ok {
				return
			}

			a.handleMessage(jw)
		}
	}
}
