package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dfl/tools/certgen/commands"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("SECRETS_ROOT_DIR", "/Users/duffleman/Source/infra-secrets/certificates")

	commands.RootCmd.AddCommand(commands.CreateCRLFileCmd)
	commands.RootCmd.AddCommand(commands.GenerateClientCertificateCmd)
	commands.RootCmd.AddCommand(commands.GenerateKeyPairCmd)
	commands.RootCmd.AddCommand(commands.GenerateRootCACmd)
	commands.RootCmd.AddCommand(commands.GenerateServerCertificateCmd)
	commands.RootCmd.AddCommand(commands.InteractiveCmd)
	commands.RootCmd.AddCommand(commands.VersionCmd)
}

func main() {
	viper.SetEnvPrefix("CERTGEN")
	viper.AutomaticEnv()

	if err := commands.RootCmd.Execute(); err != nil {
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
