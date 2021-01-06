package cmd

import (
	"dfl/svc/monitor/server"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor handles all the service monitoring",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := server.DefaultConfig()

		err := envconfig.Process("monitor", &cfg)
		if err != nil {
			return err
		}

		return server.Run(cfg)
	},
}
