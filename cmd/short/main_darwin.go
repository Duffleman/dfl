package main

import (
	"dfl/lib/keychain/darwin"
	authCommands "dfl/tools/auth/commands"
	"dfl/tools/short/commands"

	"github.com/urfave/cli/v2"
)

func init() {
	kc := darwin.Keychain{}

	rootCmd.Commands = []*cli.Command{
		commands.AddShortcut(kc),
		commands.CopyURL(kc),
		commands.DeleteResource(kc),
		commands.RemoveShortcut(kc),
		commands.Screenshot(kc),
		commands.SetNSFW(kc),
		commands.ShortenURL(kc),
		commands.UploadSigned(kc),
		commands.ViewDetails(kc),

		authCommands.Login(clientID, "short:upload short:delete", kc),
	}
}
