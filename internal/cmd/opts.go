package cmd

import (
	"flag"

	"github.com/lucasepe/t2i/internal/image/color"
	"github.com/lucasepe/t2i/internal/image/text"
	cmdutil "github.com/lucasepe/t2i/internal/util/cmd"
)

type Options struct {
	text.RenderOptions
	TabSize int
	Outfile string
}

func Configure(fs *flag.FlagSet, vals *FlagValues, args []string) Options {
	// Parse command-line args
	if err := fs.Parse(args); err != nil {
		return Options{}
	}

	// Validate & parse colors
	textCol, err := color.ParseHexColor(*vals.TextColor)
	cmdutil.CheckErr("Invalid text color", err)

	bgCol, err := color.ParseHexColor(*vals.BgColor)
	cmdutil.CheckErr("Invalid background color", err)

	// Build the final render options
	return Options{
		RenderOptions: text.RenderOptions{
			ImageWidth:      *vals.Width,
			ImageHeight:     *vals.Height,
			Margin:          *vals.Margin,
			FontSize:        *vals.FontSize,
			DPI:             *vals.DPI,
			AutoSize:        *vals.AutoSize,
			Square:          *vals.Square,
			TextColor:       textCol,
			BackgroundColor: bgCol,
			LineSpacing:     *vals.LineSpacing,
		},
		TabSize: *vals.TabSize,
		Outfile: *vals.Outfile,
	}
}
