package color

import (
	"errors"
	"image/color"
	"strconv"
	"strings"
)

// ParseHexColor parses hex colors in formats:
// #RGB, RGB, #RRGGBB, RRGGBB, #RRGGBBAA, RRGGBBAA
// Returns: *image.Uniform
func ParseHexColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")

	var r, g, b, a uint8 = 0, 0, 0, 255 // default alpha = 255

	switch len(s) {
	case 3:
		// "abc" → "aabbcc"
		r16, err := strconv.ParseUint(string([]byte{s[0], s[0]}), 16, 8)
		if err != nil {
			return nil, err
		}
		g16, err := strconv.ParseUint(string([]byte{s[1], s[1]}), 16, 8)
		if err != nil {
			return nil, err
		}
		b16, err := strconv.ParseUint(string([]byte{s[2], s[2]}), 16, 8)
		if err != nil {
			return nil, err
		}
		r, g, b = uint8(r16), uint8(g16), uint8(b16)

	case 6:
		// RRGGBB
		val, err := strconv.ParseUint(s, 16, 24)
		if err != nil {
			return nil, err
		}
		r = uint8(val >> 16)
		g = uint8(val >> 8)
		b = uint8(val)

	case 8:
		// RRGGBBAA
		val, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return nil, err
		}
		r = uint8(val >> 24)
		g = uint8(val >> 16)
		b = uint8(val >> 8)
		a = uint8(val)

	default:
		return nil, errors.New("invalid hex color length")
	}

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}
