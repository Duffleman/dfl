package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"dfl/lib/cher"
	"dfl/lib/keychain"
	"dfl/svc/short"

	"github.com/spf13/cobra"
)

func ViewDetails(kc keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "view [query]",
		Aliases: []string{"v"},
		Short:   "View details of a resource",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 0 {
				return nil
			}

			return cher.New("missing_arguments", nil)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			query, err := handleQueryInput(args)
			if err != nil {
				return err
			}

			res, err := makeClient(kc).ViewDetails(ctx, &short.IdentifyResource{
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
