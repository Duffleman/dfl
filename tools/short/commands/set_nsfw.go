package commands

import (
	"context"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func SetNSFW(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "nsfw",
		ArgsUsage: "[query]",
		Aliases:   []string{"n"},
		Usage:     "Toggle the NSFW flag",

		Action: func(c *cli.Context) error {
			ctx := context.Background()

			startTime := time.Now()

			query, err := handleQueryInput(c.Args().Slice())
			if err != nil {
				return err
			}

			newState, err := toggleNSFW(ctx, kc, query)
			if err != nil {
				return err
			}

			log.Infof("NSFW flag is now %s", newState)

			log.Infof("Done in %s", time.Now().Sub(startTime))

			return nil
		},
	}
}

func toggleNSFW(ctx context.Context, kc keychain.Keychain, query string) (string, error) {
	client, err := newClient(kc)
	if err != nil {
		return "", err
	}

	res, err := client.ViewDetails(ctx, &short.IdentifyResource{
		Query: query,
	})
	if err != nil {
		return "", err
	}

	newState := "ON"

	if res.NSFW {
		newState = "OFF"
	}

	return newState, client.SetNSFW(ctx, &short.SetNSFWRequest{
		IdentifyResource: short.IdentifyResource{
			Query: query,
		},
		NSFW: !res.NSFW,
	})
}
