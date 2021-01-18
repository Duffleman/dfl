package commands

import (
	"dfl/tools/certgen/app"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var GenerateServerCertificateCmd = &cli.Command{
	Name:      "generate_server_ceritificate",
	ArgsUsage: "[domain]",
	Aliases:   []string{"gsc"},
	Usage:     "Generate a server certificate",

	Action: func(c *cli.Context) error {
		name := c.Args().First()

		rootDirectory := viper.GetString("SECRETS_ROOT_DIR")

		app := &app.App{
			RootDirectory: rootDirectory,
		}

		return app.GenerateServerCertificate(name)
	},
}
