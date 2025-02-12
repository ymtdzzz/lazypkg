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
			input: "base-files/noble-updates 13ubuntu10.2 amd64 [13ubuntu10.1 からアップグレード可]",
			want: &PackageInfo{
				Name:       "base-files",
				OldVersion: "13ubuntu10.1",
				NewVersion: "13ubuntu10.2",
			},
		},
		{
			input: "bpftrace/noble-updates 0.20.2-1ubuntu4.3 amd64 [0.20.2-1ubuntu4.2 からアップグレード可]",
			want: &PackageInfo{
				Name:       "bpftrace",
				OldVersion: "0.20.2-1ubuntu4.2",
				NewVersion: "0.20.2-1ubuntu4.3",
			},
		},
		{
			input: "linux-firmware/noble-updates 20240318.git3b128b60-0ubuntu2.9 amd64 [20240318.git3b128b60-0ubuntu2.7 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-firmware",
				OldVersion: "20240318.git3b128b60-0ubuntu2.7",
				NewVersion: "20240318.git3b128b60-0ubuntu2.9",
			},
		},
		{
			input: "linux-generic-hwe-24.04/noble-updates 6.11.0-17.17~24.04.2+2 amd64 [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-generic-hwe-24.04",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.11.0-17.17~24.04.2+2",
			},
		},
		{
			input: "linux-headers-generic-hwe-24.04/noble-updates 6.11.0-17.17~24.04.2+2 amd64 [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-headers-generic-hwe-24.04",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.11.0-17.17~24.04.2+2",
			},
		},
		{
			input: "linux-image-generic-hwe-24.04/noble-updates 6.11.0-17.17~24.04.2+2 amd64 [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-image-generic-hwe-24.04",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.11.0-17.17~24.04.2+2",
			},
		},
		{
			input: "linux-libc-dev/noble-updates 6.8.0-53.55 amd64 [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-libc-dev",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.8.0-53.55",
			},
		},
		{
			input: "linux-modules-nvidia-550-generic-hwe-24.04/noble-updates 6.11.0-17.17~24.04.2+1 amd64 [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-modules-nvidia-550-generic-hwe-24.04",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.11.0-17.17~24.04.2+1",
			},
		},
		{
			input: "linux-modules-nvidia-550-open-generic-hwe-24.04/noble-updates 6.11.0-17.17~24.04.2+1 amd64 [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-modules-nvidia-550-open-generic-hwe-24.04",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.11.0-17.17~24.04.2+1",
			},
		},
		{
			input: "linux-tools-common/noble-updates 6.8.0-53.55 all [6.8.0-52.53 からアップグレード可]",
			want: &PackageInfo{
				Name:       "linux-tools-common",
				OldVersion: "6.8.0-52.53",
				NewVersion: "6.8.0-53.55",
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
