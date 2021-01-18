package commands

import (
	"dfl/tools/certgen/app"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var GenerateRootCACmd = &cli.Command{
	Name:    "generate_root_ca",
	Aliases: []string{"gca"},
	Usage:   "Generate a new root CA",

	Action: func(c *cli.Context) error {
		rootDirectory := viper.GetString("SECRETS_ROOT_DIR")

		app := &app.App{
			RootDirectory: rootDirectory,
		}

		return app.GenerateRootCA()
	},
}
