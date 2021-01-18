package commands

import (
	"fmt"

	"dfl/tools/certgen"

	"github.com/urfave/cli/v2"
)

var RootCmd = &cli.App{
	Name:  "certgen",
	Usage: "certgen manages and generates certificates for you.",

	Commands: []*cli.Command{
		VersionCmd,
	},
}

var VersionCmd = &cli.Command{
	Name:  "version",
	Usage: "Print the version number of certgen",

	Action: func(c *cli.Context) error {
		fmt.Printf("certgen v%s\n", certgen.Version)

		return nil
	},
}
