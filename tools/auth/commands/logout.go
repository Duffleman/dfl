package commands

import (
	"fmt"

	"dfl/lib/cli"
	"dfl/lib/keychain"

	"github.com/spf13/cobra"
)

func Logout(keychain keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout and delete auth credentials",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := keychain.DeleteItem("Auth"); err != nil {
				return err
			}

			fmt.Println(cli.Success("Logged out!"))

			return nil
		},
	}
}
