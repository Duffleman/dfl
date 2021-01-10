package cli

import (
	"github.com/ttacon/chalk"
)

var Success = chalk.Green.NewStyle().
	Style

var Warning = chalk.Yellow.NewStyle().
	Style

var Danger = chalk.White.NewStyle().
	WithBackground(chalk.Red).
	Style
