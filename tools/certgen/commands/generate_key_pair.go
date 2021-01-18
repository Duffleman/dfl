package commands

import (
	"dfl/tools/certgen/app"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var GenerateKeyPairCmd = &cli.Command{
	Name:      "generate_key_pair",
	ArgsUsage: "[name]",
	Aliases:   []string{"gkp"},
	Usage:     "Generate a public and private key pair",

	Action: func(c *cli.Context) error {
		name := c.Args().First()

		rootDirectory := viper.GetString("SECRETS_ROOT_DIR")

		app := &app.App{
			RootDirectory: rootDirectory,
		}

		if err := app.GenerateKeyPair(name); err != nil {
			return err
		}

		return app.VerifyKeyPair(name)
	},
}
