package commands

import (
	"fmt"

	"dfl/lib/keychain"

	"github.com/spf13/cobra"
)

func WhoAmI(keychain keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "whoami",
		Aliases: []string{"?"},
		Short:   "Let's find out who you are",

		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(keychain)
			if err != nil {
				return err
			}

			whoami, err := client.WhoAmI(cmd.Context())
			if err != nil {
				return err
			}

			fmt.Println("User ID:  ", whoami.UserID)
			fmt.Println("Username: ", whoami.Username)

			return nil
		},
	}
}
