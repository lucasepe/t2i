package color

import (
	"image/color"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		in       string
		expected color.RGBA
	}{
		// 3-digit
		{"#fff", color.RGBA{255, 255, 255, 255}},
		{"abc", color.RGBA{0xaa, 0xbb, 0xcc, 255}},
		{"#000", color.RGBA{0, 0, 0, 255}},

		// 6-digit
		{"ffffff", color.RGBA{255, 255, 255, 255}},
		{"#112233", color.RGBA{0x11, 0x22, 0x33, 255}},
		{"000000", color.RGBA{0, 0, 0, 255}},

		// 8-digit RGBA
		{"ff000080", color.RGBA{255, 0, 0, 128}},
		{"#00ff00cc", color.RGBA{0, 255, 0, 204}},
	}

	for _, tt := range tests {
		got, err := ParseHexColor(tt.in)
		if err != nil {
			t.Errorf("unexpected error for %s: %v", tt.in, err)
			continue
		}

		if got != tt.expected {
			t.Errorf("for %s expected %v, got %v", tt.in, tt.expected, got)
		}
	}
}

func TestParseHexColorInvalid(t *testing.T) {
	invalid := []string{
		"",
		"1",
		"12",
		"12345",
		"xyz",
		"#12G", // invalid hex char
	}

	for _, s := range invalid {
		_, err := ParseHexColor(s)
		if err == nil {
			t.Errorf("expected error for invalid input %s", s)
		}
	}
}
