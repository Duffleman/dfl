package cmd

import (
	"os"

	"dfl/lib/config"
	"dfl/svc/auth/server"

	"github.com/spf13/cobra"
)

// RootCmd is the default command for the upload proxy service executable.
var RootCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth handles all auth related matters",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := server.DefaultConfig()

		err := config.FromEnvironment(os.Getenv, &cfg)
		if err != nil {
			return err
		}

		return server.Run(cfg)
	},
}
