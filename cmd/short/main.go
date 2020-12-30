package main

import (
	"fmt"
	"os"

	authCommands "dfl/tools/auth/commands"
	"dfl/tools/short/commands"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(commands.AddShortcutCmd)
	rootCmd.AddCommand(commands.CopyURLCmd)
	rootCmd.AddCommand(commands.DeleteResourceCmd)
	rootCmd.AddCommand(commands.RemoveShortcutCmd)
	rootCmd.AddCommand(commands.ScreenshotCmd)
	rootCmd.AddCommand(commands.SetNSFWCmd)
	rootCmd.AddCommand(commands.ShortenURLCmd)
	rootCmd.AddCommand(commands.UploadSignedCmd)
	rootCmd.AddCommand(commands.ViewDetailsCmd)

	rootCmd.AddCommand(authCommands.Login("client_000000C3NCrPNP0CxPAK3M1uMjeTY", "short:upload short:delete"))
}

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFL")
	viper.SetDefault("AUTH_URL", "https://auth.dfl.mn")
	viper.SetDefault("SHORT_URL", "https://dfl.mn")
	viper.SetDefault("FS", "/Users/duffleman/.dfl")

	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "short",
	Short: "CLI tool to upload images to a short server",
	Long:  "A CLI tool to manage files and URLs being uploaded and removed from your chosen short server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
