package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	clilib "dfl/lib/cli"
	"dfl/lib/keychain"
	authCommands "dfl/tools/auth/commands"
	"dfl/tools/short/app"
	"dfl/tools/short/commands"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const clientID = "client_000000C3NCrPNP0CxPAK3M1uMjeTY"

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "https://auth.dfl.mn")
	viper.SetDefault("SHORT_URL", "https://dfl.mn")

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
		Name:  "short",
		Usage: "CLI tool to upload images to a short server",

		Commands: []*cli.Command{
			commands.AddShortcut,
			commands.CopyURL,
			commands.DeleteResource,
			commands.RemoveShortcut,
			commands.Screenshot,
			commands.SetNSFW,
			commands.ShortenURL,
			commands.UploadSigned,
			commands.ViewDetails,

			authCommands.Login(clientID, "short:upload short:delete"),
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
