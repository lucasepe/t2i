package text

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/lucasepe/t2i/internal/bdf"
	"github.com/lucasepe/t2i/internal/bdf/fonts"
	"github.com/lucasepe/t2i/internal/gg"
)

// RenderOptions defines the configuration options for TextToImage.
//
// Each field controls a specific aspect of text rendering, image sizing
// and layout. Unless otherwise specified, fields default to 0 (or false)
// and are overridden by internal defaults of the renderer.
type RenderOptions struct {
	// ImageWidth is the width of the generated image in pixels.
	// If zero, the width may be determined automatically when AutoSize is true.
	ImageWidth int

	// ImageHeight is the height of the generated image in pixels.
	// If zero, the height may be determined automatically when AutoSize is true.
	ImageHeight int

	// Margin sets the margin (padding) around the text block, in pixels.
	// The margin is applied equally on all sides.
	Margin int

	// LineSpacing defines the ratio between successive text lines.
	// For example, a value of 1.3 increases line spacing by 30%.
	LineSpacing float64

	// FontSize sets the size of the font used to render text, in points.
	FontSize float64

	// DPI sets the dots-per-inch resolution for font rendering.
	DPI float64

	// TransparentBackground, when true, makes the background transparent
	// instead of solid white.
	TransparentBackground bool

	// AutoSize, when true, measures the input text and automatically
	// sets ImageWidth and ImageHeight to fit the text block.
	// Margin is still applied on top of the measured dimensions.
	AutoSize bool

	// Square, when true, forces the final image to be square-shaped,
	// with both sides equal to the larger of ImageWidth or ImageHeight.
	Square bool

	TextColor color.Color

	BackgroundColor color.Color
}

func RenderGG(text string, opts RenderOptions) (*gg.Context, error) {
	fnt, err := bdf.LoadFont(fonts.Cozette())
	if err != nil {
		return nil, fmt.Errorf("unable to load fonts: %s", err)
	}

	face := fnt.NewFace()

	// Crea un contesto provvisorio per misurare il testo
	dc := gg.NewContext(100, 100)
	dc.SetFontFace(face)

	lines := strings.Split(text, "\n")

	// Misura la larghezza massima e l'altezza totale
	maxW := 0.0
	for _, line := range lines {
		w, _ := dc.MeasureString(line)
		if w > maxW {
			maxW = w
		}
	}
	metrics := dc.FontMetrics()
	ascent := float64(metrics.Ascent.Ceil())
	descent := float64(metrics.Descent.Ceil())
	lineHeight := ascent + descent

	totalH := float64(len(lines))*lineHeight*opts.LineSpacing - (lineHeight * (opts.LineSpacing - 1))

	// Determina dimensioni immagine finali
	width := opts.ImageWidth
	height := opts.ImageHeight

	if opts.AutoSize {
		width = int(maxW) + 2*opts.Margin
		height = int(totalH) + 2*opts.Margin
	}

	if opts.Square {
		size := width
		if height > width {
			size = height
		}
		width = size
		height = size
	}

	// Crea contesto finale con dimensioni corrette
	dc = gg.NewContext(width, height)
	dc.SetFontFace(face)

	// Colore sfondo
	if opts.BackgroundColor != nil {
		dc.SetColor(opts.BackgroundColor)
		dc.Clear()
	}

	// Colore testo
	if opts.TextColor != nil {
		dc.SetColor(opts.TextColor)
	} else {
		dc.SetRGB(0, 0, 0)
	}

	// Disegna le linee centrando verticalmente e orizzontalmente
	xoff := float64(opts.Margin) / 2.0
	yoff := float64(opts.Margin) / 2.0

	x := xoff
	y := yoff + opts.LineSpacing*lineHeight + ascent
	for _, line := range lines {
		dc.DrawStringAnchored(line, x, y, 0, 0)
		y += lineHeight * opts.LineSpacing
	}

	return dc, nil
}
