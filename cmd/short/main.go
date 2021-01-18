package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const clientID = "client_000000C3NCrPNP0CxPAK3M1uMjeTY"

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "https://auth.dfl.mn")
	viper.SetDefault("SHORT_URL", "https://dfl.mn")

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
	Name:  "short",
	Usage: "CLI tool to upload images to a short server",
}
