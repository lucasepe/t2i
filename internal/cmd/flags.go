package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type FlagValues struct {
	Margin      *int
	TabSize     *int
	Square      *bool
	TextColor   *string
	BgColor     *string
	LineSpacing *float64
	Outfile     *string
	ShowHelp    *bool
}

// NewFlagSet creates a FlagSet with all supported CLI options.
func NewFlagSet() (*flag.FlagSet, *FlagValues) {
	fs := flag.NewFlagSet("t2i", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	vals := &FlagValues{
		Margin:      fs.Int("m", 24, "Margin around the text (pixels)"),
		TabSize:     fs.Int("t", 4, "Number of spaces to replace each tab"),
		Square:      fs.Bool("s", false, "Force the image to be square"),
		TextColor:   fs.String("c", "#000", "Text color in HEX format"),
		BgColor:     fs.String("b", "#fff", "Background color in HEX format"),
		LineSpacing: fs.Float64("l", 1.3, "Space between lines."),
		Outfile:     fs.String("o", "out.png", "Output PNG file path"),
		ShowHelp:    fs.Bool("h", false, "Show help"),
	}

	return fs, vals
}

func PrintFlags(fs *flag.FlagSet) {
	maxNameLen := 0
	fs.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
	})

	indentColumn := 2 + 1 + maxNameLen + 2 // 2 spaces + "-" + name + 2 spaces padding

	w := fs.Output()
	fs.VisitAll(func(f *flag.Flag) {
		usage := strings.Split(f.Usage, "\n")
		tot := len(usage)
		for i, line := range usage {
			if i == 0 {
				padding := strings.Repeat(" ", maxNameLen-len(f.Name)+2)
				fmt.Fprintf(w, "  -%s%s%s\n", f.Name, padding, line)
			} else {
				fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", indentColumn), line)
			}
			if f.DefValue != "" && (i == (tot - 1)) {
				fmt.Fprintf(w, "%s ↳ (default: %s)\n",
					strings.Repeat(" ", indentColumn), f.DefValue)
			}
		}
		fmt.Fprintln(w)
	})
}
