package app

import (
	"math/rand"
	"sync"
	"time"

	"dfl/svc/monitor"
)

var (
	// MaxRand is the maximum time to wait before starting a job
	MaxRand = 10
	// MinRand is the minimum time to wait before starting a job
	MinRand = 2
)

type jobWrap struct {
	Job     monitor.Job
	Outcome int
}

// StartWorkers starts all workers, 1 per job
func (a *App) StartWorkers() (*sync.WaitGroup, error) {
	wg := &sync.WaitGroup{}

	jobCh := make(chan jobWrap)

	go a.messageHandlerWorker(jobCh)

	for i, job := range a.Jobs {
		a.Logger.Infof("starting worker %d/%d", i+1, len(a.Jobs))
		wg.Add(1)
		go a.startWorker(wg, jobCh, job)
		wait := time.Duration(rand.Intn((MaxRand - MinRand) + MinRand))
		a.Logger.Infof("waiting for %d seconds", wait)
		time.Sleep(wait * time.Second)
	}

	return wg, nil
}

func (a *App) startWorker(wg *sync.WaitGroup, jobCh chan jobWrap, job *monitor.Job) {
	defer wg.Done()

	var outcome int

	for {
		switch job.Type {
		case "icmp":
			outcome = a.doICMP(job)
		case "https":
			outcome = a.doHTTPS(job, true)
		case "https-novalidate":
			outcome = a.doHTTPS(job, false)
		case "http":
			outcome = a.doHTTP(job, true)
		case "tcp":
			outcome = a.doTCP(job)
		default:
			a.Logger.Warnf("job type not implemented %s", job.Type)
			return
		}

		jobCh <- jobWrap{*job, outcome}

		time.Sleep(job.Interval * time.Second)
	}
}
