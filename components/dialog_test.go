package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		input         string
		maxLineLength int
		wantStr       string
		wantInt       int
	}{
		{
			input:         "aaa",
			maxLineLength: 5,
			wantStr:       "aaa",
			wantInt:       3,
		},
		{
			input:         "aaa bbb ccc",
			maxLineLength: 5,
			wantStr:       "aaa\nbbb\nccc",
			wantInt:       3,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotStr, gotInt := wrapText(tt.input, tt.maxLineLength)
			assert.Equal(t, tt.wantStr, gotStr)
			assert.Equal(t, tt.wantInt, gotInt)
		})
	}
}
