package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const clientID = "client_000000C3N8sN2HPqTVeqfOTsnjBJI"

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "https://auth.dfl.mn")

	viper.AutomaticEnv()

	if err := rootCmd.Run(os.Args); err != nil {
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

var rootCmd = &cli.App{
	Name:  "auth",
	Usage: "Manage your authentication to DFL services",
}
