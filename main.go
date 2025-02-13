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
		dryRun bool
	)

	rootCmd := &cobra.Command{
		Use:     "lazypkg",
		Short:   "A TUI package management application across package managers",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := components.NewAppModel(components.Config{
				DryRun: dryRun,
			})
			if err != nil {
				return err
			}
			defer m.Close()
			_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
			return err
		},
	}

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Perform update commands with --dry-run option")

	if err := rootCmd.Execute(); err != nil {
		log.Println("Error running program:", err)
		os.Exit(1)
	}
}
