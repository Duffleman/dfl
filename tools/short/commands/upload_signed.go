package commands

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"dfl/lib/keychain"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/koyachi/go-nude"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var ignoredFiles = []string{
	".DS_Store",
}

func UploadSigned(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:      "signed-upload",
		ArgsUsage: "[file]",
		Aliases:   []string{"u"},
		Usage:     "Upload a file to a signed URL",

		Action: func(c *cli.Context) error {
			mutex := sync.Mutex{}
			g, gctx := errgroup.WithContext(c.Context)
			startTime := time.Now()

			localFile, err := handleLocalFileInput(c.Args().Slice())
			if err != nil {
				return err
			}

			filePaths, err := scanDirectory(localFile)
			if err != nil {
				return err
			}

			if len(filePaths) == 0 {
				return cher.New("no_files", nil)
			}

			all := []string{}

			singleFile := len(filePaths) == 1

			for _, fn := range filePaths {
				filename := fn

				g.Go(func() error {
					log.Infof("Handling file: %s", filename)
					innerStart := time.Now()

					isNude, err := nude.IsNude(filename)
					if err != nil {
						return err
					}

					if isNude {
						log.Infof("Nudity detected in %s", filename)
					}

					file, err := ioutil.ReadFile(filename)
					if err != nil {
						return err
					}

					filePrepStart := time.Now()

					resource, err := prepareUpload(c.Context, kc, filename, file)
					if err != nil {
						return err
					}

					log.Infof("File prepared: %s (%s)", resource.URL, time.Now().Sub(filePrepStart))

					if isNude {
						log.Infof("Marking file as NSFW (%s)", resource.Hash)

						g.Go(func() error {
							_, err := toggleNSFW(gctx, kc, resource.Hash)
							return err
						})
					}

					mutex.Lock()
					all = append(all, resource.Hash)
					mutex.Unlock()

					if singleFile {
						writeClipboard(resource.URL)
						notify("File prepared", resource.URL)
					}

					err = sendFileAWS(resource.SignedLink, file)
					if err != nil {
						return err
					}

					if singleFile {
						notify("File uploaded", resource.URL)
					} else {
						log.Infof("File uploaded: %s", resource.URL)
					}

					log.Infof("File handled in %s", time.Now().Sub(innerStart))

					return nil
				})
			}

			if err := g.Wait(); err != nil {
				return err
			}

			if !singleFile {
				jointURL := fmt.Sprintf("%s%s", rootURL(), strings.Join(all, ","))
				log.Infof("Download TAR at: %s", jointURL)
				writeClipboard(jointURL)
			}

			log.Infof("Done in %s", time.Now().Sub(startTime))

			return nil
		},
	}
}

func prepareUpload(ctx context.Context, kc keychain.Keychain, filename string, file []byte) (*short.CreateSignedURLResponse, error) {
	contentType := http.DetectContentType(file)

	var name *string

	if filename != "" {
		_, tmpName := filepath.Split(filename)
		name = &tmpName
	}

	client, err := newClient(kc)
	if err != nil {
		return nil, err
	}

	return client.CreateSignedURL(ctx, &short.CreateSignedURLRequest{
		ContentType: contentType,
		Name:        name,
	})
}

// SendFileAWS uploads the file to AWS
func sendFileAWS(signedURL string, file []byte) error {
	contentType := http.DetectContentType(file)

	req, err := http.NewRequest("PUT", signedURL, bytes.NewReader(file))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return err
}

func scanDirectory(rootFile string) (filePaths []string, err error) {
	root, err := os.Stat(rootFile)
	if err != nil {
		return nil, err
	}

	if !root.IsDir() {
		return []string{rootFile}, nil
	}

	err = filepath.Walk(rootFile, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, f := range ignoredFiles {
			if strings.HasSuffix(path, f) {
				return nil
			}
		}

		filePaths = append(filePaths, path)

		return nil
	})

	return
}

func handleLocalFileInput(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	file, err := filePrompt.Run()
	if err != nil {
		return "", err
	}

	return file, nil
}
