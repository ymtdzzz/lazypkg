package executors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerDiffPackageFromHash(t *testing.T) {
	imageName := "image-name:latest"

	tests := []struct {
		localDigest  string
		remoteDigest string
		want         *PackageInfo
	}{
		{
			localDigest:  "ghcr.io/open-telemetry/demo@sha256:bdc9d2a52e796649d74a8c2566897d7a45441a11bc6bc68b54a5c4c06c563eb5",
			remoteDigest: "sha256:51cff8aaa53c0af334e4cd8fce3e698a3d5114dbd530f983f62c8e0c41ad3f8a",
			want: &PackageInfo{
				Name:       imageName,
				OldVersion: "bdc9d2a",
				NewVersion: "51cff8a",
			},
		},
		{
			localDigest:  "ghcr.io/open-telemetry/demo@sha256:bdc9d2a52e796649d74a8c2566897d7a45441a11bc6bc68b54a5c4c06c563eb5",
			remoteDigest: "sha256:bdc9d2a52e796649d74a8c2566897d7a45441a11bc6bc68b54a5c4c06c563eb5",
			want:         nil,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := dockerDiffPackageFromHash(imageName, tt.localDigest, tt.remoteDigest)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDockerDiffPackageFromHashErr(t *testing.T) {
	imageName := "image-name:latest"

	tests := []struct {
		localDigest  string
		remoteDigest string
	}{
		{
			localDigest:  "invalid-local-digest",
			remoteDigest: "sha256:51cff8aaa53c0af334e4cd8fce3e698a3d5114dbd530f983f62c8e0c41ad3f8a",
		},
		{
			localDigest:  "ghcr.io/open-telemetry/demo@sha256:bdc9d2a52e796649d74a8c2566897d7a45441a11bc6bc68b54a5c4c06c563eb5",
			remoteDigest: "invalid-remote-digest",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := dockerDiffPackageFromHash(imageName, tt.localDigest, tt.remoteDigest)
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	}
}
