package commands

import (
	"dfl/tools/certgen/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var GenerateClientCertificateCmd = &cli.Command{
	Name:      "generate_client_ceritificate",
	ArgsUsage: "[hostname]",
	Aliases:   []string{"gcc"},
	Usage:     "Generate a client certificate",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "password",
			Usage:    "Password for the exported certificate",
			Required: true,
		},
	},

	Action: func(c *cli.Context) error {
		name := c.Args().First()

		password := c.String("password")
		if password == "" {
			return cher.New("no_password_given", nil)
		}

		rootDirectory := viper.GetString("SECRETS_ROOT_DIR")

		app := &app.App{
			RootDirectory: rootDirectory,
		}

		return app.GenerateClientCertificate(name, password)
	},
}
