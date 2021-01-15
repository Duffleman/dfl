package commands

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"dfl/lib/keychain"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"github.com/spf13/cobra"
)

func CopyURL(kc keychain.Keychain) *cobra.Command {
	return &cobra.Command{
		Use:     "copy [url]",
		Aliases: []string{"c"},
		Short:   "Copy from a URL",
		Long:    "Copy from a URL to the dflimg server",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 0 {
				return nil
			}

			return cher.New("missing_arguments", nil)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := handleURLInput(args)
			if err != nil {
				return err
			}

			filePath, err := downloadFile(url)
			if err != nil {
				return err
			}
			defer os.Remove(*filePath)

			return UploadSigned(kc).RunE(cmd, []string{*filePath})
		},
	}
}

func downloadFile(urlStr string) (*string, error) {
	fileToCopyRes, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer fileToCopyRes.Body.Close()

	tmpName := ksuid.Generate("tmpfile").String()

	out, err := ioutil.TempFile("", tmpName)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, fileToCopyRes.Body)
	if err != nil {
		return nil, err
	}

	path := out.Name()

	return &path, nil
}

func handleURLInput(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	url, err := urlPrompt.Run()
	if err != nil {
		return "", err
	}

	return url, nil
}
