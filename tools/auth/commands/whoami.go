package commands

import (
	"fmt"

	"dfl/lib/keychain"

	"github.com/urfave/cli/v2"
)

func WhoAmI(keychain keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:    "whoami",
		Usage:   "Read the access token and find out what your ID and username is",
		Aliases: []string{"?"},

		Action: func(c *cli.Context) error {
			client, err := newClient(keychain)
			if err != nil {
				return err
			}

			whoami, err := client.WhoAmI(c.Context)
			if err != nil {
				return err
			}

			fmt.Println("User ID:  ", whoami.UserID)
			fmt.Println("Username: ", whoami.Username)

			return nil
		},
	}
}
