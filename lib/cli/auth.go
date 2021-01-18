package cli

import (
	"encoding/json"
	"fmt"

	"dfl/lib/keychain"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
)

func AuthHeader(keychain keychain.Keychain, tool string) (*string, error) {
	var authBytes []byte
	var err error

	if keychain == nil {
		return nil, nil
	}

	authBytes, err = keychain.GetItem("Auth")
	if err != nil || len(authBytes) == 0 {
		fmt.Printf("%s%s%s", "Run `", tool, " login` first!")
		return nil, err
	}

	res := auth.TokenResponse{}

	err = json.Unmarshal(authBytes, &res)
	if err != nil {
		fmt.Printf("%s%s%s", "Run `", tool, " login` again.")
		return nil, err
	}

	return ptr.String(res.AccessToken), nil
}
