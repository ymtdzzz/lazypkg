package executors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNpmPackageFromString(t *testing.T) {
	tests := []struct {
		input string
		want  *PackageInfo
	}{
		{
			input: "corepack   0.29.4  0.31.0  0.31.0  node_modules/corepack  global",
			want: &PackageInfo{
				Name:       "corepack",
				OldVersion: "0.29.4",
				NewVersion: "0.31.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := npmPackageFromString(tt.input)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNpmPackageFromStringErr(t *testing.T) {
	got, err := npmPackageFromString("Package  Current  Wanted  Latest  Location          Depended by")
	assert.Error(t, err)
	assert.Nil(t, got)
}
