package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dfl/lib/cher"
	"dfl/tools/auth/commands"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const clientID = "client_000000C3LLC0atbRyYdTsadcVM3JQ"

func init() {
	rootCmd.AddCommand(commands.Login(clientID, "auth:login"))
}

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "http://localhost:3000")
	viper.SetDefault("FS", "/Users/duffleman/.dfl")

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
