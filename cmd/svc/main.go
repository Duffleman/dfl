package main

import (
	"os"

	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/servicecontext"
	"github.com/cuvva/ksuid-go"
	"github.com/spf13/cobra"

	auth "dfl/svc/auth/server/cmd"
	monitor "dfl/svc/monitor/server/cmd"
	short "dfl/svc/short/server/cmd"
)

// RootCmd is the initial entrypoint where all services are mounted.
var RootCmd = &cobra.Command{
	Use:   "dfl",
	Short: "dfl monobinary for dfl monorepo",
	Long:  "dfl monobinary contains entrypoints for all dfl services in the monorepo",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		env := config.EnvironmentName(os.Getenv)

		ksuid.SetEnvironment(env)
		servicecontext.Set(cmd.Use, env)
	},
}

func init() {
	RootCmd.AddCommand(auth.RootCmd)
	RootCmd.AddCommand(monitor.RootCmd)
	RootCmd.AddCommand(short.RootCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
