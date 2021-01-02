package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dfl/lib/cher"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const clientID = "client_000000C3N8sN2HPqTVeqfOTsnjBJI"

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "https://auth.dfl.mn")

	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		if v, ok := err.(cher.E); ok {
			bytes, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "auth",
	Short: "CLI tool to manage auth for DFL svcs",
	Long:  "A CLI tool to manage auth things for dfl services",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
