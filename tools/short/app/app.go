package app

import (
	"fmt"
	"strings"

	"dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"
	"dfl/svc/short"

	"github.com/spf13/viper"
)

type App struct {
	RootURL    string
	AuthURL    string
	Keychain   keychain.Keychain
	Client     short.Service
	AuthClient auth.Service
}

func New(kc keychain.Keychain) (*App, error) {
	bearerToken, err := cli.AuthHeader(kc, "short")
	if err != nil {
		return nil, err
	}

	rootURL := strings.TrimSuffix(viper.GetString("SHORT_URL"), "/")
	authURL := strings.TrimSuffix(viper.GetString("AUTH_URL"), "/")

	client, err := short.NewClient(fmt.Sprintf("%s/", rootURL), bearerToken), nil
	if err != nil {
		return nil, err
	}

	authClient, err := auth.NewClient(fmt.Sprintf("%s/", authURL), bearerToken), nil
	if err != nil {
		return nil, err
	}

	return &App{
		RootURL:    rootURL,
		AuthURL:    authURL,
		Keychain:   kc,
		Client:     client,
		AuthClient: authClient,
	}, nil
}

func (a *App) GetAuthClient() auth.Service {
	return a.AuthClient
}

func (a *App) GetKeychain() keychain.Keychain {
	return a.Keychain
}
