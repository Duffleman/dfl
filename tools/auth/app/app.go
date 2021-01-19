package app

import (
	"fmt"
	"strings"

	"dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/spf13/viper"
)

type App struct {
	RootURL  string
	Keychain keychain.Keychain
	Client   auth.Service
}

func New(kc keychain.Keychain) (*App, error) {
	bearerToken, err := cli.AuthHeader(kc, "auth")
	if err != nil {
		return nil, err
	}

	rootURL := strings.TrimSuffix(viper.GetString("AUTH_URL"), "/")

	client, err := auth.NewClient(fmt.Sprintf("%s/", rootURL), bearerToken), nil
	if err != nil {
		return nil, err
	}

	return &App{
		RootURL:  rootURL,
		Keychain: kc,
		Client:   client,
	}, nil
}

func (a *App) GetAuthClient() auth.Service {
	return a.Client
}

func (a *App) GetKeychain() keychain.Keychain {
	return a.Keychain
}
