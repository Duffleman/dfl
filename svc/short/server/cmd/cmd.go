package cmd

import (
	"os"

	"dfl/lib/config"
	"dfl/svc/short/server"

	"github.com/spf13/cobra"
)

// RootCmd is the default command for the upload proxy service executable.
var RootCmd = &cobra.Command{
	Use:   "short",
	Short: "short handles all short URLs for files and URLs",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := server.DefaultConfig()

		err := config.FromEnvironment(os.Getenv, &cfg)
		if err != nil {
			return err
		}

		return server.Run(cfg)
	},
}
