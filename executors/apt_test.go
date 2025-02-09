package executors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAptPackageFromString(t *testing.T) {
	tests := []struct {
		input string
		want  *PackageInfo
	}{
		{
			input: "language-pack-en-base/now 1:24.04+20240817 all [インストール済み、1:24.04+20250130 にアップグレード可]",
			want: &PackageInfo{
				Name:    "language-pack-en-base",
				Version: "1:24.04+20240817",
				Arch:    "all",
			},
		},
		{
			input: "language-pack-en/noble-updates 1:24.04+20250130 all [1:24.04+20240817 からアップグレード可]",
			want: &PackageInfo{
				Name:    "language-pack-en",
				Version: "1:24.04+20250130",
				Arch:    "all",
			},
		},
		{
			input: "language-pack-gnome-ja-base/noble-updates 1:24.04+20250130 all [1:24.04+20240817 からアップグレード可]",
			want: &PackageInfo{
				Name:    "language-pack-gnome-ja-base",
				Version: "1:24.04+20250130",
				Arch:    "all",
			},
		},
		{
			input: "language-pack-gnome-ja/noble-updates 1:24.04+20250130 all [1:24.04+20240817 からアップグレード可]",
			want: &PackageInfo{
				Name:    "language-pack-gnome-ja",
				Version: "1:24.04+20250130",
				Arch:    "all",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := aptPackageFromString(tt.input)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAptPackageFromStringErr(t *testing.T) {
	got, err := aptPackageFromString("一覧表示... 完了")
	assert.Error(t, err)
	assert.Nil(t, got)
}
