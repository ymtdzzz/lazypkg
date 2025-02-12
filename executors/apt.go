package executors

import (
	"bufio"
	"fmt"
	"io"
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

func (ae *AptExecutor) Update(pkg, password string, dryRun bool) error {
	cmds := []string{"sudo", "-S", "apt", "install", "--only-upgrade"}
	if dryRun {
		cmds = append(cmds, "--dry-run")
	}
	cmds = append(cmds, pkg)

	log.Printf("Running %s", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = strings.NewReader(password + "\n")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var passworderr bool
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		line := scanner.Text()
		log.Print(line)
		if strings.Contains(line, "no password was provided") {
			passworderr = true
		}
	}

	if err := cmd.Wait(); err != nil {
		if passworderr {
			return PasswordErr
		}
		return err
	}

	return nil
}

func (ae *AptExecutor) BulkUpdate(pkgs []string, password string, dryRun bool) error {
	cmds := []string{"sudo", "-S", "apt", "install", "--only-upgrade"}
	if dryRun {
		cmds = append(cmds, "--dry-run")
	}
	cmds = append(cmds, pkgs...)

	log.Printf("Running %s", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = strings.NewReader(password + "\n")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var passworderr bool
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		line := scanner.Text()
		log.Print(line)
		if strings.Contains(line, "no password was provided") {
			passworderr = true
		}
	}

	if err := cmd.Wait(); err != nil {
		if passworderr {
			return PasswordErr
		}
		return err
	}

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
