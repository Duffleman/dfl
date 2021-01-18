package commands

import (
	"context"
	"strings"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteResource(kc keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "delete [query]",
		Aliases: []string{"d"},
		Short:   "Delete a resource",
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

			err = deleteResource(ctx, kc, query)
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
