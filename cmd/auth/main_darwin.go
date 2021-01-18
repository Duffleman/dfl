package main

import (
	"dfl/lib/keychain/darwin"
	"dfl/tools/auth/commands"

	"github.com/urfave/cli/v2"
)

func init() {
	kc := darwin.Keychain{}

	rootCmd.Commands = []*cli.Command{
		commands.Login(clientID, "auth:login", kc),
		commands.CreateInviteCode(kc),
		commands.Logout(kc),
		commands.Manage(kc),
		commands.Register(kc),
		commands.SetToken(kc),
		commands.ShowAccessToken(kc),
		commands.WhoAmI(kc),
	}
}
