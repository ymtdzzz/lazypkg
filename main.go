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
		dryRun         bool
		excludes       []string
		enableFeatures []string
		demo           bool
	)

	rootCmd := &cobra.Command{
		Use:     "lazypkg",
		Short:   "A TUI package management application across package managers",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := components.NewAppModel(components.NewConfig(dryRun, excludes, enableFeatures, demo))
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
	rootCmd.Flags().StringArrayVar(&enableFeatures, "enable-feature", []string{}, "Optional feature name to be enabled in lazypkg [docker]")
	rootCmd.Flags().BoolVar(&demo, "demo", false, "")

	if err := rootCmd.Flags().MarkHidden("demo"); err != nil {
		log.Println("Error marking hidden flag:", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Println("Error running program:", err)
		os.Exit(1)
	}
}
