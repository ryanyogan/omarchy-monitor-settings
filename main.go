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

// AppConfig holds application configuration
type AppConfig struct {
	NoHyprlandCheck bool
	DebugMode       bool
	ForceLiveMode   bool
	IsTestMode      bool
}

// AppServices holds all the services used by the application
type AppServices struct {
	Config          *AppConfig
	MonitorDetector MonitorDetectorInterface
	ScalingManager  ScalingManagerInterface
	ConfigManager   ConfigManagerInterface
}

// MonitorDetectorInterface defines the interface for monitor detection
type MonitorDetectorInterface interface {
	DetectMonitors() ([]Monitor, error)
}

// ScalingManagerInterface defines the interface for scaling management
type ScalingManagerInterface interface {
	GetIntelligentScalingOptions(monitor Monitor) []ScalingOption
}

// ConfigManagerInterface defines the interface for configuration management
type ConfigManagerInterface interface {
	ApplyMonitorScale(monitor Monitor, scale float64) error
	ApplyGTKScale(scale int) error
	ApplyFontDPI(dpi int) error
	ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error
}

// NewAppServices creates and configures all application services
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
		Run: func(cmd *cobra.Command, args []string) {
			config := &AppConfig{
				NoHyprlandCheck: noHyprlandCheck,
				DebugMode:       debugMode,
				ForceLiveMode:   forceLiveMode,
				IsTestMode:      false, // Normal execution
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
	// Don't start TUI during tests
	if config.IsTestMode {
		return fmt.Errorf("TUI disabled during tests")
	}

	// Create services with dependency injection
	services := NewAppServices(config)

	// Initialize the TUI model with injected services
	model := NewModelWithServices(services)

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
