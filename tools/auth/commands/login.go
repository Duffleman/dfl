package commands

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os/exec"
	"path"
	"runtime"

	"dfl/lib/cher"
	"dfl/svc/auth"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tjarratt/babble"
)

func Login(clientID, scope string) *cobra.Command {
	return &cobra.Command{
		Use:     "login",
		Aliases: []string{"l"},
		Short:   "Login",

		RunE: func(cmd *cobra.Command, args []string) error {
			original, hashed := makeCodeChallenge()
			state := makeState()

			_, err := loadFromFile("auth.json")
			if err != nil {
				return err
			}

			params := url.Values{
				"client_id":             []string{clientID},
				"scope":                 []string{scope},
				"response_type":         []string{"code"},
				"state":                 []string{state},
				"nonce":                 []string{makeNonce()},
				"code_challenge_method": []string{"S256"},
				"code_challenge":        []string{hashed},
			}

			rootURL := viper.GetString("AUTH_URL")

			url := fmt.Sprintf("%s/authorize?%s", rootURL, params.Encode())

			cw := color.New(color.BgYellow)
			ca := color.New(color.BgHiGreen)
			c := color.New()

			cw.Print("Careful")
			c.Printf(": %s: ", "Don't forget to look for the state! It should match")
			ca.Println(state)

			err = openBrowser(url)
			if err != nil {
				cw.Print("Warning")
				c.Printf(": %s\n\n", "Cannot open your browser for you, type in this URL:")
				c.Add(color.FgCyan)
				c.Printf("%s\n\n", url)
			}

			authToken, err := authTokenPrompt.Run()
			if err != nil {
				return err
			}

			res, err := makeClient().Token(context.Background(), &auth.TokenRequest{
				ClientID:     clientID,
				GrantType:    "authorization_code",
				Code:         authToken,
				CodeVerifier: original,
			})
			if err != nil {
				return err
			}

			authBytes, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(path.Join(getRootPath(), "auth.json"), authBytes, 0644)
			if err != nil {
				return err
			}

			ca.Println("Logged in!")

			return nil
		},
	}
}

var authTokenPrompt = promptui.Prompt{
	Label: "Auth token",
	Validate: func(in string) error {
		if len(in) == 0 {
			return cher.New("too_short", nil)
		}

		return nil
	},
}

func makeState() string {
	babbler := babble.NewBabbler()
	babbler.Count = 4

	return babbler.Babble()
}

func makeCodeChallenge() (original, hashed string) {
	randomBytes, err := generateRandomBytes(32)
	if err != nil {
		panic(err)
	}

	original = base64url.Encode(randomBytes)

	h := sha256.New()
	h.Write([]byte(original))
	hashed = base64url.Encode(h.Sum(nil))

	return
}

func makeNonce() string {
	bytes, err := generateRandomBytes(32)
	if err != nil {
		panic(err)
	}

	return base64url.Encode(bytes)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
