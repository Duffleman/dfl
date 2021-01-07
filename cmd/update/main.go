package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"dfl/lib/cher"

	"github.com/spf13/cobra"
)

const repoOwner = "Duffleman"
const repoName = "dfl"

var prefixes = map[string]string{
	"darwin.amd64":  "mac64",
	"windows.amd64": "win64",
}

var rootCmd = &cobra.Command{
	Use:   "update",
	Short: "CLI tool to manage CLI updates",
	Long:  "A CLI tool to manage updating other DFL CLI tools",

	RunE: func(cmd *cobra.Command, args []string) error {
		prefix := fmt.Sprintf("%s.%s", runtime.GOOS, runtime.GOARCH)

		binPrefix, ok := prefixes[prefix]
		if !ok {
			return cher.New(cher.NoLongerSupported, nil)
		}

		fmt.Println("👀 Looking for the latest release")

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return cher.New("invalid_request", nil)
		}

		var ghr GitHubRelease

		if err := json.NewDecoder(res.Body).Decode(&ghr); err != nil {
			return err
		}

		fmt.Println("💯 Found:", ghr.Name)

		var assetsForOS []Asset

		for _, asset := range ghr.Assets {
			if strings.HasPrefix(asset.Name, binPrefix) && asset.State == "uploaded" {
				assetsForOS = append(assetsForOS, asset)
			}
		}

		updateBin, err := os.Executable()
		if err != nil {
			log.Println(err)
		}

		binPath, _ := path.Split(updateBin)

		fmt.Println("📲", len(assetsForOS), "assets to download and install")

		if err := downloadAssets(assetsForOS); err != nil {
			return err
		}

		if err := moveAssets(binPrefix, binPath, assetsForOS); err != nil {
			return err
		}

		if err := cleanupAssets(assetsForOS); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if v, ok := err.(cher.E); ok {
			bytes, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(bytes))
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}
