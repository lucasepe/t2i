package cmd

import (
	_ "embed"
	"flag"
	"fmt"
)

var (
	//go:embed assets/help.txt
	helpText string

	//go:embed assets/support.txt
	supportText string
)

func Usage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Fprintln(fs.Output(), helpText)
		fmt.Fprint(fs.Output(), "FLAGS:\n\n")
		PrintFlags(fs)
		fmt.Fprintln(fs.Output(), supportText)
	}
}
