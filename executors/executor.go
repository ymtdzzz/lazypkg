package executors

import (
	"errors"
	"os/exec"
)

var PasswordErr = errors.New("password is required")

type PackageInfo struct {
	Name    string
	Version string
	Arch    string
}

type Executor interface {
	GetPackages() ([]*PackageInfo, error)
	Update(pkg, password string, dryRun bool) error
	BulkUpdate(pkgs []string, password string, dryRun bool) error
	Valid() bool
}

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
