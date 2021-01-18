package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dfl/tools/certgen/commands"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func init() {
	viper.SetDefault("SECRETS_ROOT_DIR", "/Users/duffleman/Source/infra-secrets/certificates")

	commands.RootCmd.Commands = []*cli.Command{
		commands.CreateCRLFileCmd,
		commands.GenerateClientCertificateCmd,
		commands.GenerateKeyPairCmd,
		commands.GenerateRootCACmd,
		commands.GenerateServerCertificateCmd,
		commands.InteractiveCmd,
		commands.VersionCmd,
	}
}

func main() {
	viper.SetEnvPrefix("CERTGEN")
	viper.AutomaticEnv()

	if err := commands.RootCmd.Run(os.Args); err != nil {
		if c, ok := err.(cher.E); ok {
			bytes, err := json.MarshalIndent(c, "", "  ")
			if err != nil {
				panic(err)
			}

			fmt.Println(string(bytes))
			os.Exit(1)
		}

		fmt.Println(err)
		os.Exit(1)
	}
}
