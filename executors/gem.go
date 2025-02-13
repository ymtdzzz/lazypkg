package executors

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var gemPattern = regexp.MustCompile(`^([^\s]+)\s+\(([^\s]+)\s<\s([^\s]+)\)`)

type GemExecutor struct{}

func (ge *GemExecutor) Valid() bool {
	return cmdExists("gem")
}

func (ge *GemExecutor) GetPackages(_ string) ([]*PackageInfo, error) {
	var packages []*PackageInfo

	log.Print("Running gem outdated")
	cmd := exec.Command("gem", "outdated")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// NOTE: invalid row will be skipped
		// TODO: log for verbose
		if pkg, err := gemPackageFromString(line); err == nil {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func (ge *GemExecutor) Update(pkg, _ string, dryRun bool) error {
	cmds := []string{"gem", "update", pkg}
	if dryRun {
		log.Printf("[dry-run] %s", strings.Join(cmds, " "))
		return nil
	}

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

func (ge *GemExecutor) BulkUpdate(pkgs []string, _ string, dryRun bool) error {
	cmds := []string{"gem", "update"}
	cmds = append(cmds, pkgs...)
	if dryRun {
		log.Printf("[dry-run] %s", strings.Join(cmds, " "))
		return nil
	}

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

func (ge *GemExecutor) Close() {}

func gemPackageFromString(input string) (*PackageInfo, error) {
	matches := gemPattern.FindStringSubmatch(input)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid input provided: %s", input)
	}
	return &PackageInfo{
		Name:       matches[1],
		OldVersion: matches[2],
		NewVersion: matches[3],
	}, nil
}
