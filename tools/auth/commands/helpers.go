package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"dfl/svc/auth"

	"github.com/spf13/viper"
)

func makeClient() auth.Service {
	return auth.NewClient(rootURL())
}

func rootURL() string {
	return fmt.Sprintf("%s/", viper.Get("AUTH_URL").(string))
}

func getRootPath() string {
	return viper.Get("FS").(string)
}

func loadFromFile(filename string) (res *auth.TokenResponse, err error) {
	path := path.Join(getRootPath(), filename)

	file, err := os.OpenFile(path, os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	json.NewDecoder(file).Decode(&res)

	return
}
