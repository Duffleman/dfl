package commands

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"dfl/lib/keychain"
	"dfl/svc/auth"
)

var ShowAccessTokenCmd = &cobra.Command{
	Use:     "show-access-token",
	Aliases: []string{"sat"},
	Short:   "Show the currently stored access token",

	RunE: func(cmd *cobra.Command, args []string) error {
		var authBytes []byte
		var err error

		if !keychain.Supported() {
			authBytes, err = readFromFile(path.Join(getRootPath(), "auth.json"))
		} else {
			authBytes, err = keychain.GetItem("Auth")
		}

		if err != nil {
			return err
		}

		var res auth.TokenResponse

		if err := json.Unmarshal(authBytes, &res); err != nil {
			return err
		}

		fmt.Println(res.AccessToken)

		return nil
	},
}
