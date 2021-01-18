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

func SetNSFW(kc keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "nsfw [query]",
		Aliases: []string{"n"},
		Short:   "Toggle the NSFW flag",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 0 {
				return nil
			}

			return cher.New("missing_arguments", nil)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			startTime := time.Now()

			query, err := handleQueryInput(args)
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
