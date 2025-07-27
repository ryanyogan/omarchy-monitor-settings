package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ryanyogan/omarchy-monitor-settings/internal/app"
	"github.com/ryanyogan/omarchy-monitor-settings/internal/tui"
	"github.com/spf13/cobra"
)

var (
	noHyprlandCheck bool
	debugMode       bool
	forceLiveMode   bool
	version         = "dev"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "omarchy-monitor-settings",
		Short:   "A stunning TUI for managing monitor resolution and scaling",
		Long:    "A beautiful terminal interface for detecting and configuring monitor resolution, scaling, and font settings in Hyprland/Wayland environments.",
		Version: version,
		Run: func(_ *cobra.Command, _ []string) {
			config := &app.Config{
				NoHyprlandCheck: noHyprlandCheck,
				DebugMode:       debugMode,
				ForceLiveMode:   forceLiveMode,
				IsTestMode:      false,
			}

			if err := runTUI(config); err != nil {
				log.Fatalf("Error running TUI: %v", err)
			}
		},
	}

	rootCmd.Flags().BoolVar(&noHyprlandCheck, "no-hyprland-check", false, "Skip Hyprland environment check (useful for testing)")
	rootCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode")
	rootCmd.Flags().BoolVar(&forceLiveMode, "force-live", false, "Force live mode (bypass all checks for testing)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI(config *app.Config) error {
	if config.IsTestMode {
		return fmt.Errorf("TUI disabled during tests")
	}

	// Check if we're in a test environment by looking for test flags
	if testing.Testing() {
		return fmt.Errorf("TUI disabled during tests")
	}

	services := app.NewServices(config)

	model := tui.NewModelWithServices(services)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
