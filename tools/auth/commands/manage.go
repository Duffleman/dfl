package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"dfl/lib/keychain"
)

func Manage(keychain keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "manage",
		Aliases: []string{"m"},
		Short:   "Manage credentials online",

		RunE: func(cmd *cobra.Command, args []string) error {
			return openBrowser(fmt.Sprintf("%s/u2f_manage", strings.TrimSuffix(rootURL(), "/")))
		},
	}
}
