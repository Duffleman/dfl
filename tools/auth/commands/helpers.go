package commands

import (
	"fmt"
	"io/ioutil"
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

func writeToFile(in string, data []byte) error {
	dir, _ := path.Split(in)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(in, data, 0700); err != nil {
		return err
	}

	return nil
}

func readFromFile(in string) (data []byte, err error) {
	return ioutil.ReadFile(in)
}
