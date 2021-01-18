package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"dfl/lib/keychain"

	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"github.com/urfave/cli/v2"
)

const screenshotCmd = "screencapture -i"
const timeout = 1 * time.Minute

func Screenshot(kc keychain.Keychain) *cli.Command {
	return &cli.Command{
		Name:  "screenshot",
		Usage: "Take a screenshot & upload it",

		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithTimeout(c.Context, timeout)
			defer cancel()

			tmpName := fmt.Sprintf("%s-*.png", ksuid.Generate("file").String())
			out, err := ioutil.TempFile("", tmpName)
			if err != nil {
				return err
			}
			defer out.Close()

			err = exec.CommandContext(ctx, "screencapture", "-i", out.Name()).Run()
			if err != nil {
				return err
			}
			defer os.Remove(out.Name())

			tmpFile, err := os.Stat(out.Name())
			if os.IsNotExist(err) {
				return nil
			}

			if tmpFile.Size() == 0 {
				notify("Cancelled", "No image was captured.")
				return nil
			}

			return UploadSigned(kc).Action(c)
		},
	}
}
