package executors

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var homebrewPattern = regexp.MustCompile(`(\S+) \(([^)]+)\) < (\S+)`)

type HomebrewExecutor struct{}

func (he *HomebrewExecutor) Valid() bool {
	return cmdExists("brew")
}

func (he *HomebrewExecutor) GetPackages(_ string) ([]*PackageInfo, error) {
	var packages []*PackageInfo

	// check for update
	log.Print("Running brew update")
	cmd := exec.Command("brew", "update")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	log.Print("Running brew outdated --verbose")
	cmd = exec.Command("brew", "outdated", "--verbose")
	output, err = cmd.Output()
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

func (he *HomebrewExecutor) Update(pkg, _ string, dryRun bool) error {
	cmds := []string{"brew", "upgrade"}
	if dryRun {
		cmds = append(cmds, "--dry-run")
	}
	cmds = append(cmds, pkg)

	log.Printf("Running %s", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		log.Print(line)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (he *HomebrewExecutor) BulkUpdate(pkgs []string, password string, dryRun bool) error {
	cmds := []string{"brew", "upgrade"}
	if dryRun {
		cmds = append(cmds, "--dry-run")
	}
	cmds = append(cmds, pkgs...)

	log.Printf("Running %s", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		log.Print(line)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (he *HomebrewExecutor) Close() {}

func homebrewPackageFromString(input string) (*PackageInfo, error) {
	matches := homebrewPattern.FindStringSubmatch(input)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid input provided: %s", input)
	}
	return &PackageInfo{
		Name:       matches[1],
		OldVersion: matches[2],
		NewVersion: matches[3],
	}, nil
}
