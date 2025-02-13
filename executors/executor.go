package executors

import (
	"errors"
	"os/exec"
)

var ErrPassword = errors.New("password is required")

type PackageInfo struct {
	Name       string
	OldVersion string
	NewVersion string
}

type Executor interface {
	GetPackages(password string) ([]*PackageInfo, error)
	Update(pkg, password string, dryRun bool) error
	BulkUpdate(pkgs []string, password string, dryRun bool) error
	Valid() bool
	Close()
}

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
