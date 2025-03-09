package executors

import (
	"errors"
	"os/exec"
)

// ErrPassword is returned when a required password is not provided
var ErrPassword = errors.New("password is required")

// PackageInfo represents information about a package including its name and version details
type PackageInfo struct {
	Name       string
	OldVersion string
	NewVersion string
}

// Executor defines the interface for package management operations
type Executor interface {
	// GetPackages retrieves a list of available package updates.
	// The password parameter is required for package managers that need elevated privileges.
	GetPackages(password string) ([]*PackageInfo, error)

	// Update performs an update operation on a single package.
	// If dryRun is true, it will only simulate the update without making actual changes.
	// The password parameter is required for package managers that need elevated privileges.
	Update(pkg, password string, dryRun bool) error

	// BulkUpdate performs update operations on multiple packages simultaneously.
	// If dryRun is true, it will only simulate the updates without making actual changes.
	// The password parameter is required for package managers that need elevated privileges.
	BulkUpdate(pkgs []string, password string, dryRun bool) error

	// Valid checks if the package manager is available and usable on the current system.
	Valid() bool

	// Close performs any necessary cleanup operations when the executor is no longer needed.
	Close()
}

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
