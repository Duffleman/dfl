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

func RemoveShortcut(kc keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "remove-shortcut [query] [shortcut]",
		Aliases: []string{"rsc"},
		Short:   "Remove a shortcut",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 || len(args) == 0 {
				return nil
			}

			return cher.New("missing_arguments", nil)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			startTime := time.Now()

			query, shortcut, err := handleShortcutInput(args)
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
