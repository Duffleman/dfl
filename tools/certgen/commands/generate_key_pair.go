package commands

import (
	"dfl/tools/certgen/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GenerateKeyPairCmd = &cobra.Command{
	Use:     "generate_key_pair [name]",
	Aliases: []string{"gkp"},
	Short:   "Generate a public and private key pair",
	Args:    cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		rootDirectory := viper.GetString("SECERTS_ROOT_DIR")

		app := &app.App{
			RootDirectory: rootDirectory,
		}

		if err := app.GenerateKeyPair(name); err != nil {
			return err
		}

		return app.VerifyKeyPair(name)
	},
}
