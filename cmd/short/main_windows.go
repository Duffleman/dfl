package main

import (
	"dfl/lib/keychain/windows"
	authCommands "dfl/tools/auth/commands"
	"dfl/tools/short/commands"
)

func init() {
	kc := windows.Keychain{}

	rootCmd.AddCommand(commands.AddShortcut(kc))
	rootCmd.AddCommand(commands.CopyURL(kc))
	rootCmd.AddCommand(commands.DeleteResource(kc))
	rootCmd.AddCommand(commands.RemoveShortcut(kc))
	rootCmd.AddCommand(commands.Screenshot(kc))
	rootCmd.AddCommand(commands.SetNSFW(kc))
	rootCmd.AddCommand(commands.ShortenURL(kc))
	rootCmd.AddCommand(commands.UploadSigned(kc))
	rootCmd.AddCommand(commands.ViewDetails(kc))

	rootCmd.AddCommand(authCommands.Login("client_000000C3NCrPNP0CxPAK3M1uMjeTY", "short:upload short:delete", kc))
}