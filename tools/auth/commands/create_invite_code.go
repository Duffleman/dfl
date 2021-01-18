package commands

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"dfl/lib/cli"
	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var scopesRegex = regexp.MustCompile(`^(?:[a-z*]+:[a-z*]+\s?)+$`)

func CreateInviteCode(kc keychain.Keychain) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-invite-code",
		Aliases: []string{"cic"},
		Short:   "Create an invite code for another user to register",
		Args:    cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := context.Background()

			scopes, code, expiresAt, err := handleInputs()
			if err != nil {
				return err
			}

			client, err := newClient(kc)
			if err != nil {
				return err
			}

			res, err := client.CreateInviteCode(ctx, &auth.CreateInviteCodeRequest{
				Code:      code,
				ExpiresAt: expiresAt,
				Scopes:    scopes,
			})
			if err != nil {
				return err
			}

			fmt.Println(cli.Success("Success!"))

			fmt.Printf("Code: %s\n", res.Code)
			if res.ExpiresAt != nil {
				fmt.Printf("Expires at: %s\n", res.ExpiresAt.Format(time.RFC3339))
			}

			return nil
		},
	}

	return cmd
}

func handleInputs() (scopes string, code *string, expiresAt *time.Time, err error) {
	scopes, err = scopesPrompt.Run()
	if err != nil {
		return "", nil, nil, err
	}

	codeStr, err := codePrompt.Run()
	if err != nil {
		return "", nil, nil, err
	}

	expiresAtStr, err := expiresAtPrompt.Run()
	if err != nil {
		return "", nil, nil, err
	}

	if codeStr != "" {
		code = &codeStr
	}

	if expiresAtStr != "" {
		expiresAtTS, _ := time.Parse(time.RFC3339, expiresAtStr)
		expiresAt = &expiresAtTS
	}

	return
}

var scopesPrompt = promptui.Prompt{
	Label: "Scopes",
	Validate: func(in string) error {
		if len(in) == 0 {
			return cher.New("too_short", nil)
		}

		if !scopesRegex.MatchString(in) {
			return cher.New("invalid_scope_format", nil)
		}

		return nil
	},
}

var codePrompt = promptui.Prompt{
	Label: "Code",
	Validate: func(in string) error {
		return nil
	},
}

var expiresAtPrompt = promptui.Prompt{
	Label: "Expires at",
	Validate: func(in string) error {
		if len(in) == 0 {
			return nil
		}

		_, err := time.Parse(time.RFC3339, in)
		return err
	},
}
