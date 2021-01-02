package commands

import (
	"encoding/json"
	"fmt"

	"dfl/lib/keychain"
	"dfl/svc/auth"
	"dfl/svc/short"

	"github.com/atotto/clipboard"
	b "github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// AppName for notifications
const AppName = "DFL Short"

func notify(title, body string) {
	err := b.Notify(title, body, "")
	if err != nil {
		log.Warn(err)
	}
}

func makeClient(keychain keychain.Keychain) short.Service {
	return short.NewClient(rootURL(), getAuthHeader(keychain))
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("SHORT_URL").(string))
}

func getAuthHeader(keychain keychain.Keychain) string {
	var authBytes []byte
	var err error

	authBytes, err = keychain.GetItem("Auth")
	if err != nil {
		fmt.Println("Run `short login` first!")
		panic(err)
	}

	res := auth.TokenResponse{}

	err = json.Unmarshal(authBytes, &res)
	if err != nil {
		fmt.Println("Run `short login` again.")
		panic(err)
	}

	return fmt.Sprintf("%s %s", res.TokenType, res.AccessToken)
}

func writeClipboard(in string) {
	err := clipboard.WriteAll(in)
	if err != nil {
		log.Warn("Could not copy to clipboard.")
	}
}
