package commands

import (
	"dfl/lib/keychain"

	"github.com/fatih/color"
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

			c := color.New(color.BgGreen)
			c.Println("Logged out!")

			return nil
		},
	}
}
