package executors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomebrewPackageFromString(t *testing.T) {
	tests := []struct {
		input string
		want  *PackageInfo
	}{
		{
			input: "fastfetch (2.33.0) < 2.35.0",
			want: &PackageInfo{
				Name:       "fastfetch",
				OldVersion: "2.33.0",
				NewVersion: "2.35.0",
			},
		},
		{
			input: "openjdk (23.0.1) < 23.0.2",
			want: &PackageInfo{
				Name:       "openjdk",
				OldVersion: "23.0.1",
				NewVersion: "23.0.2",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := homebrewPackageFromString(tt.input)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHomebrewPackageFromStringErr(t *testing.T) {
	got, err := homebrewPackageFromString("==> Downloading https://formulae.brew.sh/api/formula.jws.json")
	assert.Error(t, err)
	assert.Nil(t, got)
}
