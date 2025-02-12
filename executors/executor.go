package executors

import "errors"

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
}
