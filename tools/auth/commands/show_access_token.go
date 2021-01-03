package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"dfl/lib/cher"
	dfljwt "dfl/lib/jwt"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func ShowAccessToken(keychain keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "show-access-token",
		Aliases: []string{"sat"},
		Short:   "Show the currently stored access token",

		RunE: func(cmd *cobra.Command, args []string) error {
			var authBytes []byte
			var err error

			authBytes, err = keychain.GetItem("Auth")
			if err != nil {
				return err
			}

			var res auth.TokenResponse
			var dflclaims dfljwt.DFLClaims

			if err := json.Unmarshal(authBytes, &res); err != nil {
				return err
			}

			if token, _ := jwt.ParseWithClaims(res.AccessToken, &dflclaims, nil); token == nil {
				return cher.New("cannot_parse_token", nil)
			}

			fmt.Fprintf(os.Stdout, res.AccessToken)

			fmt.Fprintf(os.Stderr, "\n\n")
			fmt.Fprintf(os.Stderr, "User ID:    %s\n", dflclaims.Subject)
			fmt.Fprintf(os.Stderr, "Username:   %s\n", dflclaims.Username)
			fmt.Fprintf(os.Stderr, "Scopes:     %s\n", dflclaims.Scope)
			fmt.Fprintf(os.Stderr, "Client ID:  %s\n", dflclaims.Audience)
			fmt.Fprintf(os.Stderr, "Issuer:     %s\n", dflclaims.Issuer)
			fmt.Fprintf(os.Stderr, "Expires at: ")

			expiresAt := time.Unix(dflclaims.ExpiresAt, 0)

			c := color.New()

			now := time.Now()

			if now.After(expiresAt) {
				c.Add(color.BgRed)
			} else {
				c.Add(color.BgGreen)
			}

			duration := expiresAt.Sub(now)

			c.Fprintf(os.Stderr, "%s", expiresAt.Format(time.RFC3339))
			fmt.Fprintf(os.Stderr, " (%s)\n", duration.Round(time.Second))

			return nil
		},
	}
}
