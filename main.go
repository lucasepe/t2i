package main

import (
	"fmt"
	"os"

	"github.com/lucasepe/t2i/internal/cmd"
	imageio "github.com/lucasepe/t2i/internal/image/io"
	"github.com/lucasepe/t2i/internal/image/text"
	cmdutil "github.com/lucasepe/t2i/internal/util/cmd"
	ioutil "github.com/lucasepe/t2i/internal/util/io"
	textutil "github.com/lucasepe/t2i/internal/util/text"
)

func main() {
	fs, fv := cmd.NewFlagSet()

	// Custom usage
	fs.Usage = cmd.Usage(fs)

	opts := cmd.Configure(fs, fv, os.Args[1:])

	// If help flag was supplied: print help and exit
	if *fv.ShowHelp {
		fs.Usage()
		return
	}

	input, err := ioutil.ReadInput(fs.Args())
	cmdutil.CheckErr("error reading input", err)
	input = textutil.Clean(input, opts.TabSize)

	img, err := text.Render(input, opts.RenderOptions)
	cmdutil.CheckErr("rendering text to image failed", err)

	f, err := os.Create(opts.Outfile)
	cmdutil.CheckErr("cannot write output", err)
	defer f.Close()

	err = imageio.WriteToFile(img, opts.Outfile, imageio.PNG)
	cmdutil.CheckErr("failed to write image", err)

	fmt.Println("Image written to", opts.Outfile)
}
