package executors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGemPackageFromString(t *testing.T) {
	tests := []struct {
		input string
		want  *PackageInfo
	}{
		{
			input: "mini_magick (4.13.2 < 5.1.2)",
			want: &PackageInfo{
				Name:       "mini_magick",
				OldVersion: "4.13.2",
				NewVersion: "5.1.2",
			},
		},
		{
			input: "opentelemetry-instrumentation-base (0.22.6 < 0.23.0)",
			want: &PackageInfo{
				Name:       "opentelemetry-instrumentation-base",
				OldVersion: "0.22.6",
				NewVersion: "0.23.0",
			},
		},
		{
			input: "strscan (3.0.9 < 3.1.2)",
			want: &PackageInfo{
				Name:       "strscan",
				OldVersion: "3.0.9",
				NewVersion: "3.1.2",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := gemPackageFromString(tt.input)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGemPackageFromStringErr(t *testing.T) {
	got, err := gemPackageFromString("invalid input")
	assert.Error(t, err)
	assert.Nil(t, got)
}
