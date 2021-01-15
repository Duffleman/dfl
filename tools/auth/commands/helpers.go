package commands

import (
	"fmt"

	"dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/spf13/viper"
)

func makeClient(keychain keychain.Keychain) auth.Service {
	if keychain == nil {
		return auth.NewClient(rootURL(), nil)
	}

	authHeader := cli.AuthHeader(keychain, "auth")

	return auth.NewClient(rootURL(), &authHeader)
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("AUTH_URL").(string))
}
