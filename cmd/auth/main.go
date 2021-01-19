package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	clilib "dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/tools/auth/app"
	"dfl/tools/auth/commands"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const clientID = "client_000000C3N8sN2HPqTVeqfOTsnjBJI"

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "https://auth.dfl.mn")

	viper.AutomaticEnv()

	if err := rootCmd.Run(os.Args); err != nil {
		if v, ok := err.(cher.E); ok {
			bytes, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

var rootCmd *cli.App

func makeRoot(kc keychain.Keychain) {
	rootCmd = &cli.App{
		Name:  "auth",
		Usage: "Manage your authentication to DFL services",

		Commands: []*cli.Command{
			commands.CreateInviteCode,
			commands.Login(clientID, "auth:login"),
			commands.Logout,
			commands.Manage,
			commands.Register,
			commands.SetToken,
			commands.ShowAccessToken,
			commands.WhoAmI,
		},

		Before: func(c *cli.Context) error {
			app, err := app.New(kc)
			if err != nil {
				return err
			}

			c.Context = context.WithValue(c.Context, clilib.AppKey, app)

			return nil
		},
	}
}
