package commands

import (
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"dfl/lib/keychain"
)

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and delete auth credentials",

	RunE: func(cmd *cobra.Command, args []string) error {
		if !keychain.Supported() {
			if err := os.Remove(path.Join(getRootPath(), "auth.json")); err != nil {
				return err
			}
		} else {
			keychain.GetItem("Auth")

			if err := keychain.DeleteItem("Auth"); err != nil {
				return err
			}
		}

		c := color.New(color.BgGreen)
		c.Println("Logged out!")

		return nil
	},
}
