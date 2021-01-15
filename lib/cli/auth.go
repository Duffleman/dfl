package cli

import (
	"encoding/json"
	"fmt"

	"dfl/lib/keychain"
	"dfl/svc/auth"
)

func AuthHeader(keychain keychain.Keychain, tool string) string {
	if keychain == nil {
		return ""
	}

	var authBytes []byte
	var err error

	authBytes, err = keychain.GetItem("Auth")
	if err != nil {
		fmt.Printf("%s%s%s", "Run `", tool, " login` first!")
		panic(err)
	}

	res := auth.TokenResponse{}

	err = json.Unmarshal(authBytes, &res)
	if err != nil {
		fmt.Printf("%s%s%s", "Run `", tool, " login` again.")
		panic(err)
	}

	return res.AccessToken
}
