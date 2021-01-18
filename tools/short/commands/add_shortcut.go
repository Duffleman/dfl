package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func AddShortcut(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "add-shortcut",
		Usage:     "Add a shortcut",
		ArgsUsage: "[query] [shortcut]",
		Aliases:   []string{"asc"},

		Action: func(c *cli.Context) error {
			startTime := time.Now()

			query, shortcut, err := handleShortcutInput(c.Args().Slice())
			if err != nil {
				return err
			}

			err = addShortcut(c.Context, kc, query, shortcut)
			if err != nil {
				return err
			}

			writeClipboard(fmt.Sprintf("%s/:%s", rootURL(), shortcut))
			notify("Added shortcut", fmt.Sprintf("%s/:%s", rootURL(), shortcut))

			log.Infof("Done in %s", time.Now().Sub(startTime))

			return nil
		},
	}
}

func addShortcut(ctx context.Context, kc keychain.Keychain, query, shortcut string) error {
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

	return client.AddShortcut(ctx, body)
}

func handleShortcutInput(args []string) (string, string, error) {
	if len(args) == 2 {
		return strings.TrimPrefix(args[0], rootURL()), args[1], nil
	}

	query, err := queryPrompt.Run()
	if err != nil {
		return "", "", err
	}

	shortcut, err := shortcutPrompt.Run()
	if err != nil {
		return "", "", err
	}

	return strings.TrimPrefix(query, rootURL()), shortcut, nil
}
