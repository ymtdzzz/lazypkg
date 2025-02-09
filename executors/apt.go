package executors

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var aptPattern = regexp.MustCompile(`^([a-zA-Z0-9\+\-\.]+)\/([^\s]+)\s+([a-zA-Z0-9\+\-\.\:]+)\s+([a-zA-Z0-9\+\-\.]+)`)

type AptExecutor struct{}

func (ae *AptExecutor) GetPackages() ([]*PackageInfo, error) {
	var packages []*PackageInfo

	// check for update
	log.Print("Running apt list --upgradable")
	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// NOTE: invalid row will be skipped
		// TODO: log for verbose
		if pkg, err := aptPackageFromString(line); err == nil {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func (ae *AptExecutor) Update(pkg string) error {
	log.Printf("Running apt install --dry-run --only-upgrade %s", pkg)
	cmd := exec.Command("sudo", "apt", "install", "--dry-run", "--only-upgrade", pkg)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		log.Print(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (ae *AptExecutor) BulkUpdate(pkgs []string) error {
	return nil
}

func aptPackageFromString(input string) (*PackageInfo, error) {
	matches := aptPattern.FindStringSubmatch(input)
	if len(matches) < 5 {
		return nil, fmt.Errorf("invalid input provided: %s", input)
	}
	return &PackageInfo{
		Name:    matches[1],
		Version: matches[3],
		Arch:    matches[4],
	}, nil
}
