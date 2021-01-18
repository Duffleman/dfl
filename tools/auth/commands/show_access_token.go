package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	authlib "dfl/lib/auth"
	clilib "dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/cli/v2"
)

func ShowAccessToken(keychain keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:    "Show access token",
		Usage:   "Show the currently stored access token",
		Aliases: []string{"sat"},

		Action: func(c *cli.Context) error {
			var authBytes []byte
			var err error

			authBytes, err = keychain.GetItem("Auth")
			if err != nil {
				return err
			}

			var res auth.TokenResponse
			var dflclaims authlib.DFLClaims

			if err := json.Unmarshal(authBytes, &res); err != nil {
				return err
			}

			if token, _ := jwt.ParseWithClaims(res.AccessToken, &dflclaims, nil); token == nil {
				return cher.New("cannot_parse_token", nil)
			}

			fmt.Fprintf(os.Stdout, res.AccessToken)

			fmt.Fprintf(os.Stderr, "\n\n")
			fmt.Fprintf(os.Stderr, "Version:    %s\n", dflclaims.Version)
			fmt.Fprintf(os.Stderr, "User ID:    %s\n", dflclaims.Subject)
			fmt.Fprintf(os.Stderr, "Username:   %s\n", dflclaims.Username)
			fmt.Fprintf(os.Stderr, "Scopes:     %s\n", dflclaims.Scopes)
			fmt.Fprintf(os.Stderr, "Client ID:  %s\n", dflclaims.Audience)
			fmt.Fprintf(os.Stderr, "Issuer:     %s\n", dflclaims.Issuer)
			fmt.Fprintf(os.Stderr, "Expires at: ")

			expiresAt := time.Unix(dflclaims.ExpiresAt, 0)

			now := time.Now()
			duration := expiresAt.Sub(now)

			var style func(string) string

			if now.After(expiresAt) {
				style = clilib.Danger
			} else {
				style = clilib.Success
			}

			fmt.Fprintf(os.Stderr, style(expiresAt.Format(time.RFC3339)))
			fmt.Fprintf(os.Stderr, " (%s)\n", duration.Round(time.Second))

			return nil
		},
	}
}
