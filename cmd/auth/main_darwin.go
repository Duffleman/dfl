package main

import (
	"dfl/lib/keychain/darwin"
	"dfl/tools/auth/commands"
)

func init() {
	kc := darwin.Keychain{}

	rootCmd.AddCommand(commands.Login(clientID, "auth:login", kc))
	rootCmd.AddCommand(commands.Logout(kc))
	rootCmd.AddCommand(commands.CreateInviteCode(kc))
	rootCmd.AddCommand(commands.Manage(kc))
	rootCmd.AddCommand(commands.Register(kc))
	rootCmd.AddCommand(commands.ShowAccessToken(kc))
}
