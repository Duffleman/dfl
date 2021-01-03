package main

import (
	"dfl/lib/keychain/windows"
	"dfl/tools/auth/commands"
)

func init() {
	kc := windows.Keychain{}

	rootCmd.AddCommand(commands.Login(clientID, "auth:login", kc))
	rootCmd.AddCommand(commands.Logout(kc))
	rootCmd.AddCommand(commands.Manage(kc))
	rootCmd.AddCommand(commands.ShowAccessToken(kc))
}
