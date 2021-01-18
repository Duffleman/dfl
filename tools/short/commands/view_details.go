package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"dfl/lib/keychain"
	"dfl/svc/short"

	"github.com/urfave/cli/v2"
)

func ViewDetails(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "view",
		ArgsUsage: "[query]",
		Aliases:   []string{"v"},
		Usage:     "View details of a resource",

		Action: func(c *cli.Context) error {
			ctx := context.Background()

			query, err := handleQueryInput(c.Args().Slice())
			if err != nil {
				return err
			}

			client, err := newClient(kc)
			if err != nil {
				return err
			}

			res, err := client.ViewDetails(ctx, &short.IdentifyResource{
				Query: query,
			})
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}
}
