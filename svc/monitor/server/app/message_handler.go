package app

import (
	sdk "github.com/andygrunwald/cachet"
)

func (a *App) handleMessage(jw jobWrap) {
	if jw.Job.ComponentID == 0 {
		a.Logger.Warnf("cannot update component %s, no matched catchet component", jw.Job.Name)
		return
	}

	_, _, err := a.Cachet.Components.Update(jw.Job.ComponentID, &sdk.Component{
		Status: jw.Outcome,
	})
	if err != nil {
		a.Logger.WithError(err).Warnf("cannot update component %s", jw.Job.ComponentName)
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
