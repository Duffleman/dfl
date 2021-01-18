package commands

import (
	"crypto/x509"

	"dfl/tools/certgen/app"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var filename = "crl.txt"

var crlTemplate = &x509.RevocationList{
	SignatureAlgorithm: x509.ECDSAWithSHA384,
}

var CreateCRLFileCmd = &cli.Command{
	Name:    "create_crl_file",
	Aliases: []string{"crl"},
	Usage:   "Create a CRL file",

	Action: func(c *cli.Context) error {
		rootDirectory := viper.GetString("SECRETS_ROOT_DIR")

		app := &app.App{
			RootDirectory: rootDirectory,
		}

		return app.GenerateCRL()
	},
}
