package commands

import (
	authlib "dfl/lib/auth"
	"dfl/lib/keychain"
	"dfl/svc/auth"
	"encoding/json"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

func SetToken(keychain keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:   "set [token]",
		Short: "Override the token to a new one you manually provide.",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			var dflclaims authlib.DFLClaims

			if token, _ := jwt.ParseWithClaims(args[0], &dflclaims, nil); token == nil {
				return cher.New("cannot_parse_token", nil)
			}

			res := auth.TokenResponse{
				UserID:      dflclaims.Subject,
				AccessToken: args[0],
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
