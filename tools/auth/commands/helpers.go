package commands

import (
	"encoding/json"
	"fmt"

	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/spf13/viper"
)

func makeClient(keychain keychain.Keychain) auth.Service {
	return auth.NewClient(rootURL(), getAuthHeader(keychain))
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("AUTH_URL").(string))
}

func getAuthHeader(keychain keychain.Keychain) string {
	var authBytes []byte
	var err error

	authBytes, err = keychain.GetItem("Auth")
	if err != nil {
		fmt.Println("Run `auth login` first!")
		panic(err)
	}

	res := auth.TokenResponse{}

	err = json.Unmarshal(authBytes, &res)
	if err != nil {
		fmt.Println("Run `auth login` again.")
		panic(err)
	}

	return fmt.Sprintf("%s %s", res.TokenType, res.AccessToken)
}
