package commands

import (
	"context"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ShortenURL(kc keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "shorten [url]",
		Aliases: []string{"s"},
		Short:   "Shorten a URL",
		Long:    "Shorten a URL",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 0 {
				return nil
			}

			return cher.New("missing_arguments", nil)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			startTime := time.Now()

			url, err := handleURLInput(args)
			if err != nil {
				return err
			}

			body, err := makeClient(kc).ShortenURL(ctx, &short.CreateURLRequest{
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
