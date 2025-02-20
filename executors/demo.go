package executors

import (
	"log"
	"time"
)

type DemoExecutor struct {
	cmd  string
	pkgs []*PackageInfo
}

func NewDemoExecutor(cmd string, pkgs []*PackageInfo) *DemoExecutor {
	return &DemoExecutor{
		pkgs: pkgs,
	}
}

func (de *DemoExecutor) Valid() bool {
	return true
}

func (de *DemoExecutor) GetPackages(_ string) ([]*PackageInfo, error) {
	time.Sleep(500 * time.Millisecond)
	return de.pkgs, nil
}

func (de *DemoExecutor) Update(pkg, _ string, _ bool) error {
	return de.update(pkg)
}

func (de *DemoExecutor) BulkUpdate(pkgs []string, _ string, _ bool) error {
	for _, pkg := range pkgs {
		de.update(pkg)
	}

	return nil
}

func (de *DemoExecutor) Close() {}

func (de *DemoExecutor) update(pkg string) error {
	// simulate updating a package
	for i, p := range de.pkgs {
		if p.Name == pkg {
			log.Printf("[Demo] Start to update %s", pkg)
			log.Printf("[Demo] Running %s command to update %s package ...", de.cmd, pkg)

			time.Sleep(500 * time.Millisecond)

			log.Print("Updating ...")

			time.Sleep(500 * time.Millisecond)

			log.Print("Update completed")

			// update complete and delete pkg from outdated package list
			de.pkgs = append(de.pkgs[:i], de.pkgs[i+1:]...)
		}
	}

	return nil
}
