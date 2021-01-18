package commands

import (
	"context"
	"strings"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func DeleteResource(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		ArgsUsage: "[query]",
		Aliases:   []string{"d"},
		Usage:     "Delete a resource",

		Action: func(c *cli.Context) error {
			startTime := time.Now()

			query, err := handleQueryInput(c.Args().Slice())
			if err != nil {
				return err
			}

			err = deleteResource(c.Context, kc, query)
			if err != nil {
				return err
			}

			notify("Resource deleted", query)

			log.Infof("Done in %s", time.Now().Sub(startTime))

			return nil
		},
	}
}

func deleteResource(ctx context.Context, kc keychain.Keychain, urlStr string) error {
	body := &short.IdentifyResource{
		Query: urlStr,
	}

	client, err := newClient(kc)
	if err != nil {
		return err
	}

	return client.DeleteResource(ctx, body)
}

func handleQueryInput(args []string) (string, error) {
	if len(args) == 1 {
		return strings.TrimPrefix(args[0], rootURL()), nil
	}

	query, err := queryPrompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(query, rootURL()), nil
}
