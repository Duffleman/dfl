package commands

import (
	"fmt"
	"strings"

	"dfl/lib/keychain"

	"github.com/urfave/cli/v2"
)

func Manage(keychain keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:    "manage",
		Usage:   "Manage your U2F tokens",
		Aliases: []string{"m"},

		Action: func(c *cli.Context) error {
			return openBrowser(fmt.Sprintf("%s/u2f_manage", strings.TrimSuffix(rootURL(), "/")))
		},
	}
}

func Register(keychain keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:    "register",
		Usage:   "Register for an account",
		Aliases: []string{"r"},

		Action: func(c *cli.Context) error {
			return openBrowser(fmt.Sprintf("%s/register", strings.TrimSuffix(rootURL(), "/")))
		},
	}
}
