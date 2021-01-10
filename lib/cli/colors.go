package cli

import (
	"github.com/fatih/color"
)

var Success = func(in string) string {
	c := color.New(color.BgGreen)

	return c.Sprintf(in)
}

var Warning = func(in string) string {
	c := color.New(color.FgYellow)

	return c.Sprintf(in)
}

var Danger = func(in string) string {
	c := color.New(color.BgRed)

	return c.Sprintf(in)
}
