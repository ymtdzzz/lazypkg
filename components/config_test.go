package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	type input struct {
		dryRun   bool
		excludes []string
		enables  []string
		demo     bool
	}

	tests := []struct {
		input input
		want  Config
	}{
		{
			input: input{false, []string{}, []string{}, false},
			want: Config{
				DryRun:         false,
				Excludes:       map[string]bool{},
				EnableFeatures: map[string]bool{},
			},
		},
		{
			input: input{false, []string{"hoge", "fuga"}, []string{"piyo"}, false},
			want: Config{
				DryRun: false,
				Excludes: map[string]bool{
					"fuga": true,
					"hoge": true,
				},
				EnableFeatures: map[string]bool{
					"piyo": true,
				},
			},
		},
	}

	for _, tt := range tests {
		got := NewConfig(tt.input.dryRun, tt.input.excludes, tt.input.enables, tt.input.demo)
		assert.Equal(t, tt.want, got)
	}
}
