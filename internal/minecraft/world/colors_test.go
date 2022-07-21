package world

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{"#34c832", color.RGBA{0x34, 0xc8, 0x32, 0xff}},
		{"#34c83212", color.RGBA{0x34, 0xc8, 0x32, 0x12}},
		{"#abc", color.RGBA{0xaa, 0xbb, 0xcc, 0xff}},
		{"#bcda", color.RGBA{0xbb, 0xcc, 0xdd, 0xaa}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := mustParseColor(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
