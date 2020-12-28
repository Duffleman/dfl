package main

import (
	"github.com/spf13/cobra"

	auth "dfl/svc/auth/server/cmd"
)

// RootCmd is the initial entrypoint where all services are mounted.
var RootCmd = &cobra.Command{
	Use:   "dfl",
	Short: "dfl monobinary for dfl monorepo",
	Long:  "dfl monobinary contains entrypoints for all dfl services in the monorepo",
}

func init() {
	RootCmd.AddCommand(auth.RootCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
