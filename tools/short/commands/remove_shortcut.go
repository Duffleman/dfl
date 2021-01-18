package commands

import (
	"context"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func RemoveShortcut(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "remove-shortcut",
		ArgsUsage: "[query] [shortcut]",
		Aliases:   []string{"rsc"},
		Usage:     "Remove a shortcut",

		Action: func(c *cli.Context) error {
			ctx := context.Background()

			startTime := time.Now()

			query, shortcut, err := handleShortcutInput(c.Args().Slice())
			if err != nil {
				return err
			}

			err = removeShortcut(ctx, kc, query, shortcut)
			if err != nil {
				return err
			}

			notify("Removed shortcut", shortcut)

			log.Infof("Done in %s", time.Now().Sub(startTime))

			return nil
		},
	}
}

func removeShortcut(ctx context.Context, kc keychain.Keychain, query, shortcut string) error {
	body := &short.ChangeShortcutRequest{
		IdentifyResource: short.IdentifyResource{
			Query: query,
		},
		Shortcut: shortcut,
	}

	client, err := newClient(kc)
	if err != nil {
		return err
	}

	return client.RemoveShortcut(ctx, body)
}
