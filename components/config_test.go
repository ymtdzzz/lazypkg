package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	type input struct {
		dryRun   bool
		excludes []string
	}

	tests := []struct {
		input input
		want  Config
	}{
		{
			input: input{false, []string{}},
			want: Config{
				DryRun:   false,
				Excludes: map[string]bool{},
			},
		},
		{
			input: input{false, []string{"hoge", "fuga"}},
			want: Config{
				DryRun: false,
				Excludes: map[string]bool{
					"fuga": true,
					"hoge": true,
				},
			},
		},
	}

	for _, tt := range tests {
		got := NewConfig(tt.input.dryRun, tt.input.excludes)
		assert.Equal(t, tt.want, got)
	}
}
