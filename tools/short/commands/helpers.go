package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

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

func makeClient() short.Service {
	return short.NewClient(rootURL(), getAuthHeader())
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("SHORT_URL").(string))
}

func getRootPath() string {
	return viper.Get("FS").(string)
}

func getAuthHeader() string {
	path := path.Join(getRootPath(), "auth.json")

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Run `short login` first!")
		panic(err)
	}

	res := auth.TokenResponse{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		fmt.Println("Run `short login` first!")
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
