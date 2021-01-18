package commands

import (
	"fmt"

	clilib "dfl/lib/cli"
	"dfl/lib/keychain"

	"github.com/urfave/cli/v2"
)

func Logout(keychain keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:  "logout",
		Usage: "Delete your local credentials",

		Action: func(c *cli.Context) error {
			if err := keychain.DeleteItem("Auth"); err != nil {
				return err
			}

			fmt.Println(clilib.Success("Logged out!"))

			return nil
		},
	}
}
