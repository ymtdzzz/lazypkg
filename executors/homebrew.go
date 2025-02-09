package executors

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var homebrewPattern = regexp.MustCompile(`(\S+) \(([^)]+)\) < (\S+)`)

type HomebrewExecutor struct{}

func (ae *HomebrewExecutor) GetPackages() ([]*PackageInfo, error) {
	var packages []*PackageInfo

	// check for update
	log.Print("Running brew outdated --verbose")
	cmd := exec.Command("brew", "outdated", "--verbose")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// NOTE: invalid row will be skipped
		// TODO: log for verbose
		if pkg, err := homebrewPackageFromString(line); err == nil {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func (ae *HomebrewExecutor) Update(pkg string) error {
	return nil
}

func (ae *HomebrewExecutor) BulkUpdate(pkgs []string) error {
	return nil
}

func homebrewPackageFromString(input string) (*PackageInfo, error) {
	matches := homebrewPattern.FindStringSubmatch(input)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid input provided: %s", input)
	}
	return &PackageInfo{
		Name:    matches[1],
		Version: matches[2],
		Arch:    "",
	}, nil
}
