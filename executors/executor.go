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
	Update(pkg, password string) error
	BulkUpdate(pkgs []string, password string) error
}
