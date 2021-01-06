package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"dfl/lib/cher"
	"dfl/svc/monitor"
	"dfl/svc/monitor/server/app"
	"dfl/svc/monitor/server/lib/cachet"
	jobslib "dfl/svc/monitor/server/lib/jobs"

	cachetSDK "github.com/andygrunwald/cachet"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger *logrus.Logger

	Debug bool `envconfig:"debug"`

	CachetURL string `envconfig:"cachet_url"`
	CachetKey string `envconfig:"cachet_key"`

	JobsFile *string `envcofnig:"jobs_file"`
	Jobs     []byte  `envconfig:"jobs"`
}

func DefaultConfig() Config {
	return Config{
		Logger: logrus.New(),

		Debug: true,

		CachetURL: "https://status.dfl.mn",
		CachetKey: "",

		Jobs: []byte(`[{"name":"google","component_name":"Google","type":"https","host":"google.com","interval":5}]`),
	}
}

func Run(cfg Config) error {
	cfg.Logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
	}

	cachetClient, err := cachetSDK.NewClient(cfg.CachetURL, nil)
	if err != nil {
		return fmt.Errorf("cannot make cachet client: %w", err)
	}

	_, _, err = cachetClient.General.Ping()
	if err != nil {
		return fmt.Errorf("cannot ping cachet: %w", err)
	}

	cachetClient.Authentication.SetTokenAuth(cfg.CachetKey)

	jobs, err := loadJobs(cfg.Jobs, cfg.JobsFile)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
		},
	}

	clientNoValidate := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	app := &app.App{
		Logger:           cfg.Logger,
		Cachet:           &cachet.Client{Client: cachetClient},
		Client:           client,
		ClientNoValidate: clientNoValidate,
		Jobs:             jobs,
	}

	if err := app.SyncWithCachet(); err != nil {
		return fmt.Errorf("cannot sync with cachet: %w", err)
	}

	if cfg.Debug {
		cfg.Logger.Info("setting debug ON")
	} else {
		cfg.Logger.Info("setting debug OFF")
		cfg.Logger.SetLevel(logrus.WarnLevel)
	}

	wg, err := app.StartWorkers()
	if err != nil {
		return fmt.Errorf("cannot start workers: %w", err)
	}

	defer wg.Wait()

	cfg.Logger.Infof("Server running")

	return nil
}

func loadJobs(jobs []byte, jobFile *string) ([]*monitor.Job, error) {
	if len(jobs) == 0 && jobFile == nil {
		return nil, cher.New("invalid_configuration", nil)
	}

	if jobFile == nil {
		return jobslib.ParseData(jobs)
	}

	return jobslib.ParseFromFile(*jobFile)
}
