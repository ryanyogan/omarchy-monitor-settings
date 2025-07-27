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
	forceLiveMode   bool
	version         = "dev"
)

type AppConfig struct {
	NoHyprlandCheck bool
	DebugMode       bool
	ForceLiveMode   bool
	IsTestMode      bool
}

type AppServices struct {
	Config          *AppConfig
	MonitorDetector MonitorDetectorInterface
	ScalingManager  ScalingManagerInterface
	ConfigManager   ConfigManagerInterface
}

type MonitorDetectorInterface interface {
	DetectMonitors() ([]Monitor, error)
}

type ScalingManagerInterface interface {
	GetIntelligentScalingOptions(monitor Monitor) []ScalingOption
}

type ConfigManagerInterface interface {
	ApplyMonitorScale(monitor Monitor, scale float64) error
	ApplyGTKScale(scale int) error
	ApplyFontDPI(dpi int) error
	ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error
}

func NewAppServices(config *AppConfig) *AppServices {
	return &AppServices{
		Config:          config,
		MonitorDetector: NewMonitorDetector(),
		ScalingManager:  NewScalingManager(),
		ConfigManager:   NewConfigManager(config.IsTestMode),
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "omarchy-monitor-settings",
		Short:   "A stunning TUI for managing monitor resolution and scaling",
		Long:    "A beautiful terminal interface for detecting and configuring monitor resolution, scaling, and font settings in Hyprland/Wayland environments.",
		Version: version,
		Run: func(_ *cobra.Command, _ []string) {
			config := &AppConfig{
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

func runTUI(config *AppConfig) error {
	if config.IsTestMode {
		return fmt.Errorf("TUI disabled during tests")
	}

	services := NewAppServices(config)

	model := NewModelWithServices(services)

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
