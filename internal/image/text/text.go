package text

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
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

// Render renders text from a byte slice into an image.Image, applying
// the given options such as font size, DPI, margins, spacing and auto-sizing.
//
// If AutoSize is true, the renderer measures the text block and adjusts
// the image size automatically. Margin is always applied as a minimum
// padding around the text block.
//
// If Square is true, the final image will be forced to be square with side
// length rounded up to the nearest multiple of 24. In this case, margins
// are adjusted so that the text remains visually centered.
func Render(in []byte, opts RenderOptions) (image.Image, error) {
	wri, err := newTextRenderer(opts)
	if err != nil {
		return nil, err
	}

	// 1. Calcola dimensioni testo se richiesto
	measuredW, measuredH := 0, 0
	if opts.AutoSize {
		measuredW, measuredH, err = wri.measureText(in)
		if err != nil {
			return nil, err
		}
	} else {
		measuredW = opts.ImageWidth
		measuredH = opts.ImageHeight
	}

	// 2. Aggiungi margine minimo
	wri.width = measuredW + 2*wri.margin
	wri.height = measuredH + 2*wri.margin

	// 3. Se Square, porta dimensioni a multiplo di 24 (per eccesso)
	if opts.Square {
		size := math.Max(float64(wri.width), float64(wri.height))
		size = math.Ceil(size/24.0) * 24 // multiplo di 24
		wri.width = int(size)
		wri.height = int(size)
	}

	// 4. Calcola offset per mantenere centrato il blocco di testo
	extraW := wri.width - (measuredW + 2*wri.margin)
	extraH := wri.height - (measuredH + 2*wri.margin)

	wri.offsetX = extraW / 2
	wri.offsetY = extraH / 2

	return wri.imageFromText(in)
}

func newTextRenderer(opts RenderOptions) (*textRenderer, error) {
	wri, err := defaultTextRenderer()
	if err != nil {
		return nil, err
	}

	if opts.TextColor == nil {
		wri.textColor = color.Black
	} else {
		wri.textColor = opts.TextColor
	}

	if opts.BackgroundColor == nil {
		if !opts.TransparentBackground {
			wri.backgroundColor = color.White
		} else {
			wri.backgroundColor = color.Transparent
		}
	} else {
		wri.backgroundColor = opts.BackgroundColor
	}

	if opts.DPI > 0 {
		wri.dpi = opts.DPI
	}

	if opts.ImageWidth > 0 {
		wri.width = opts.ImageWidth
	}

	if opts.ImageHeight > 0 {
		wri.height = opts.ImageHeight
	}

	if opts.Margin > 0 {
		wri.margin = opts.Margin
	}

	if opts.FontSize > 0 {
		wri.fontSize = opts.FontSize
	}

	if opts.LineSpacing > 0 {
		wri.lineSpacing = opts.LineSpacing
	}

	return wri, nil
}

func defaultTextRenderer() (tr *textRenderer, err error) {
	tr = &textRenderer{
		width:           640,
		height:          640,
		dpi:             120,
		fontSize:        12,
		lineSpacing:     1.3,
		margin:          12 * 2,
		textColor:       color.Black,
		backgroundColor: color.White,
	}

	tr.ttf, err = freetype.ParseFont(gomono.TTF)
	if err != nil {
		return nil, fmt.Errorf("unable to parse font: %s", err)
	}

	tr.face = truetype.NewFace(tr.ttf, &truetype.Options{
		Size:    tr.fontSize,
		DPI:     tr.dpi,
		Hinting: font.HintingFull,
	})

	return tr, nil
}

// textRenderer is the internal implementation of ImageWriter.
type textRenderer struct {
	width           int
	height          int
	margin          int
	lineSpacing     float64
	textColor       color.Color
	backgroundColor color.Color
	dpi             float64
	fontSize        float64
	ttf             *truetype.Font
	offsetX         int // offset extra per centratura
	offsetY         int // offset extra per centratura
	face            font.Face
}

// imageFromText renders the buffered text into an RGBA image and returns it.
func (tr *textRenderer) imageFromText(in []byte) (img *image.RGBA, err error) {
	// Foreground (text) and background colors.
	fg, bg := image.NewUniform(tr.textColor), image.NewUniform(tr.backgroundColor)
	img = image.NewRGBA(image.Rect(0, 0, tr.width, tr.height))

	// Draw background.
	draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

	// Setup freetype context for text rendering.
	ctx := freetype.NewContext()
	ctx.SetDPI(tr.dpi)
	ctx.SetFont(tr.ttf)
	ctx.SetFontSize(tr.fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(fg)
	ctx.SetHinting(font.HintingFull)

	var (
		line   = 0
		reader = bytes.NewBuffer(in)
	)

	for {
		s, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return img, err
		}

		if s != "" {
			if err := tr.drawTextLine(ctx, s, line); err != nil {
				return img, err
			}
		}

		if err == io.EOF {
			break
		}
		line++
	}

	return img, nil
}

// drawTextLine draws a single line of text on the image at the given line number.
func (tr *textRenderer) drawTextLine(ctx *freetype.Context, text string, line int) error {
	text = strings.TrimRight(text, "\n")

	// Compute line height based on font size.
	//fontHeight := int(ctx.PointToFixed(tr.fontSize) >> 6)

	metrics := tr.face.Metrics()
	ascent := int(metrics.Ascent >> 6)
	lineHeight := int((metrics.Ascent + metrics.Descent) >> 6)

	// Apply line spacing and top margin.
	x := tr.margin + tr.offsetX
	y := tr.margin + tr.offsetY + ascent + line*int(float64(lineHeight)*tr.lineSpacing)

	ctx.SetFont(tr.ttf)
	ctx.SetFontSize(tr.fontSize)

	_, err := ctx.DrawString(text, freetype.Pt(x, y))
	return err
}

// MeasureText legge il testo da un io.Reader e calcola width/height del blocco
// di testo (senza margini esterni).
// MeasureText calcola la larghezza e altezza del blocco testo (senza margini esterni).
func (tr *textRenderer) measureText(in []byte) (int, int, error) {
	// face := truetype.NewFace(tr.ttf, &truetype.Options{
	// 	Size:    tr.fontSize,
	// 	DPI:     tr.dpi,
	// 	Hinting: font.HintingFull,
	// })
	// defer func() {
	// 	if c, ok := face.(io.Closer); ok {
	// 		c.Close()
	// 	}
	// }()

	var (
		reader = bytes.NewBuffer(in)
		lines  []string
	)

	for {
		s, err := reader.ReadString('\n')
		if s != "" {
			lines = append(lines, strings.TrimRight(s, "\n"))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, err
		}
	}

	if len(lines) == 0 {
		return 0, 0, nil
	}

	metrics := tr.face.Metrics()
	ascent := float64(metrics.Ascent >> 6)
	lineHeight := float64((metrics.Ascent + metrics.Descent) >> 6)

	maxW := 0.0
	d := &font.Drawer{Face: tr.face}
	for _, ln := range lines {
		adv := d.MeasureString(ln)
		if float64(adv>>6) > maxW {
			maxW = float64(adv >> 6)
		}
	}

	totalH := ascent + float64(len(lines)-1)*lineHeight*tr.lineSpacing

	return int(maxW), int(totalH), nil
}

// MeasureString returns the rendered width and height of a single string.
/*
func (tr *textRenderer) measureString(face font.Face, s string) (w, h float64) {
	d := &font.Drawer{Face: face}
	adv := d.MeasureString(s)
	metrics := face.Metrics()
	fontHeight := float64((metrics.Ascent + metrics.Descent) >> 6)

	return float64(adv >> 6), fontHeight
}
*/
