package commands

import (
	"fmt"

	"dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/spf13/viper"
)

func newClient(keychain keychain.Keychain) (auth.Service, error) {
	bearerToken, err := cli.AuthHeader(keychain, "auth")
	if err != nil {
		return nil, err
	}

	return auth.NewClient(rootURL(), bearerToken), nil
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("AUTH_URL").(string))
}
