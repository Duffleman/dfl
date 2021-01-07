package main

import (
	"dfl/lib/keychain/darwin"
	authCommands "dfl/tools/auth/commands"
	"dfl/tools/short/commands"
)

func init() {
	kc := darwin.Keychain{}

	rootCmd.AddCommand(commands.AddShortcut(kc))
	rootCmd.AddCommand(commands.CopyURL(kc))
	rootCmd.AddCommand(commands.DeleteResource(kc))
	rootCmd.AddCommand(commands.RemoveShortcut(kc))
	rootCmd.AddCommand(commands.Screenshot(kc))
	rootCmd.AddCommand(commands.SetNSFW(kc))
	rootCmd.AddCommand(commands.ShortenURL(kc))
	rootCmd.AddCommand(commands.UploadSigned(kc))
	rootCmd.AddCommand(commands.ViewDetails(kc))

	rootCmd.AddCommand(authCommands.Login(clientID, "short:upload short:delete", kc))
}
