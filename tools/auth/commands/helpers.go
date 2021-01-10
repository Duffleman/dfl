package commands

import (
	"fmt"

	"dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/spf13/viper"
)

func makeClient(keychain keychain.Keychain) auth.Service {
	return auth.NewClient(rootURL(), cli.AuthHeader(keychain, "auth"))
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("AUTH_URL").(string))
}
