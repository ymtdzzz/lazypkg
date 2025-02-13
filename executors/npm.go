package executors

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var npmPattern = regexp.MustCompile(`^([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+global`)

type NpmExecutor struct{}

func (ne *NpmExecutor) Valid() bool {
	return cmdExists("npm")
}

func (he *NpmExecutor) GetPackages(_ string) ([]*PackageInfo, error) {
	var packages []*PackageInfo

	log.Print("Running npm outdated -g")
	cmd := exec.Command("npm", "outdated", "-g")

	// NOTE: npm outdated -g returns exit code 1 even if succeeded
	output, _ := cmd.Output()

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// NOTE: invalid row will be skipped
		// TODO: log for verbose
		if pkg, err := npmPackageFromString(line); err == nil {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func (he *NpmExecutor) Update(pkg, _ string, dryRun bool) error {
	cmds := []string{"npm", "update", "-g"}
	if dryRun {
		cmds = append(cmds, "--dry-run")
	}
	cmds = append(cmds, pkg)

	log.Printf("Running %s", strings.Join(cmds, " "))
	// #nosec G204: commands are not input values
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

func (he *NpmExecutor) BulkUpdate(pkgs []string, _ string, dryRun bool) error {
	cmds := []string{"npm", "update", "-g"}
	if dryRun {
		cmds = append(cmds, "--dry-run")
	}
	cmds = append(cmds, pkgs...)

	log.Printf("Running %s", strings.Join(cmds, " "))
	// #nosec G204: commands are not input values
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

func (he *NpmExecutor) Close() {}

func npmPackageFromString(input string) (*PackageInfo, error) {
	matches := npmPattern.FindStringSubmatch(input)
	if len(matches) < 6 {
		return nil, fmt.Errorf("invalid input provided: %s", input)
	}
	return &PackageInfo{
		Name:       matches[1],
		OldVersion: matches[2],
		NewVersion: matches[3], // Wanted version
	}, nil
}
