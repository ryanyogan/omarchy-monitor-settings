package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	noHyprlandCheck bool
	debugMode       bool
	version         = "0.1.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "hyprland-monitor-tui",
		Short:   "A stunning TUI for managing Hyprland monitor resolution and scaling",
		Long:    "A beautiful terminal interface for detecting and configuring monitor resolution, scaling, and font settings in Hyprland/Wayland environments.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runTUI(); err != nil {
				log.Fatalf("Error running TUI: %v", err)
			}
		},
	}

	rootCmd.Flags().BoolVar(&noHyprlandCheck, "no-hyprland-check", false, "Skip Hyprland environment check (useful for testing)")
	rootCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI() error {
	// Initialize the TUI model
	model := NewModel()

	// Configure bubbletea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start the TUI
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
