package commands

import (
	authlib "dfl/lib/auth"
	"dfl/lib/keychain"
	"dfl/svc/auth"
	"encoding/json"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/cli/v2"
)

func SetToken(keychain keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "set-access-token",
		Usage:     "Manually set the access token",
		ArgsUsage: "[token]",

		Action: func(c *cli.Context) error {
			var dflclaims authlib.DFLClaims

			if token, _ := jwt.ParseWithClaims(c.Args().First(), &dflclaims, nil); token == nil {
				return cher.New("cannot_parse_token", nil)
			}

			res := auth.TokenResponse{
				UserID:      dflclaims.Subject,
				AccessToken: c.Args().First(),
				TokenType:   "Bearer",
				Expires:     int(dflclaims.ExpiresAt),
			}

			authBytes, err := json.Marshal(res)
			if err != nil {
				return err
			}

			return keychain.UpsertItem("Auth", authBytes)
		},
	}
}
