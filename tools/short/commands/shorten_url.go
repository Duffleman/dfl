package commands

import (
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func ShortenURL(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "shorten",
		ArgsUsage: "[url]",
		Aliases:   []string{"s"},
		Usage:     "Shorten a URL",

		Action: func(c *cli.Context) error {
			startTime := time.Now()

			url, err := handleURLInput(c.Args().Slice())
			if err != nil {
				return err
			}

			client, err := newClient(kc)
			if err != nil {
				return err
			}

			body, err := client.ShortenURL(c.Context, &short.CreateURLRequest{
				URL: url,
			})
			if err != nil {
				return err
			}

			writeClipboard(body.URL)
			notify("URL Shortened", body.URL)

			log.Infof("Done in %s: %s", time.Now().Sub(startTime), body.URL)

			return nil
		},
	}
}
