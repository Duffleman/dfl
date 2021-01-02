package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ShowAccessTokenCmd = &cobra.Command{
	Use:     "show-access-token",
	Aliases: []string{"sat"},
	Short:   "Show the currently stored access token",

	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := loadFromFile("auth.json")
		if err != nil {
			return err
		}

		fmt.Println(auth.AccessToken)

		return nil
	},
}
