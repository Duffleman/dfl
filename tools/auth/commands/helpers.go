package commands

import (
	"fmt"

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
