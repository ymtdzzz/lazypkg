package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/ymtdzzz/lazypkg/components"
)

var version = "unknown"

func main() {
	var (
		dryRun   bool
		excludes []string
	)

	rootCmd := &cobra.Command{
		Use:     "lazypkg",
		Short:   "A TUI package management application across package managers",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := components.NewAppModel(components.NewConfig(dryRun, excludes))
			if err != nil {
				return err
			}
			defer m.Close()
			_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
			return err
		},
	}

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Perform update commands with --dry-run option")
	rootCmd.Flags().StringArrayVar(&excludes, "exclude", []string{}, "Package manager name to be excluded in lazypkg")

	if err := rootCmd.Execute(); err != nil {
		log.Println("Error running program:", err)
		os.Exit(1)
	}
}
