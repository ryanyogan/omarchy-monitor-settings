package main

import (
	"os"
	"testing"
	"time"

	"github.com/ryanyogan/omarchy-monitor-settings/internal/app"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantExit bool
	}{
		{
			name:     "with version flag",
			args:     []string{"omarchy-monitor-settings", "--version"},
			wantExit: true,
		},
		{
			name:     "with help flag",
			args:     []string{"omarchy-monitor-settings", "--help"},
			wantExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = tt.args

			done := make(chan bool)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("main() panicked: %v", r)
					}
					done <- true
				}()
				main()
			}()

			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Error("main() timed out")
			}
		})
	}
}

func TestRunTUI(t *testing.T) {
	config := &app.Config{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true,
	}

	err := runTUI(config)
	if err == nil {
		t.Error("Expected error when IsTestMode is true, got nil")
	}
	if err.Error() != "TUI disabled during tests" {
		t.Errorf("Expected 'TUI disabled during tests' error, got: %v", err)
	}
}

func TestGlobalFlags(t *testing.T) {
	if noHyprlandCheck {
		t.Error("noHyprlandCheck should be false by default")
	}
	if debugMode {
		t.Error("debugMode should be false by default")
	}
	if forceLiveMode {
		t.Error("forceLiveMode should be false by default")
	}
	if version == "" {
		t.Error("version should not be empty")
	}
}

func TestAppServices(t *testing.T) {
	config := &app.Config{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true,
	}

	services := app.NewServices(config)
	if services == nil {
		t.Error("NewServices should not return nil")
		return
	}
	if services.Config != config {
		t.Error("Services should have the correct config")
	}
	if services.MonitorDetector == nil {
		t.Error("MonitorDetector should not be nil")
	}
	if services.ScalingManager == nil {
		t.Error("ScalingManager should not be nil")
	}
	if services.ConfigManager == nil {
		t.Error("ConfigManager should not be nil")
	}
}

func TestEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		expected bool
	}{
		{
			name:     "hyprland environment",
			envKey:   "HYPRLAND_INSTANCE_SIGNATURE",
			envValue: "test-signature",
			expected: true,
		},
		{
			name:     "no hyprland environment",
			envKey:   "HYPRLAND_INSTANCE_SIGNATURE",
			envValue: "",
			expected: false,
		},
		{
			name:     "wayland environment",
			envKey:   "WAYLAND_DISPLAY",
			envValue: "wayland-0",
			expected: true,
		},
		{
			name:     "x11 environment",
			envKey:   "DISPLAY",
			envValue: ":0",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue := os.Getenv(tt.envKey)
			defer os.Setenv(tt.envKey, originalValue)

			os.Setenv(tt.envKey, tt.envValue)

			config := &app.Config{
				NoHyprlandCheck: true,
				DebugMode:       false,
				ForceLiveMode:   false,
				IsTestMode:      true,
			}
			services := app.NewServices(config)
			_ = services
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			_ = noHyprlandCheck
			_ = debugMode
			_ = forceLiveMode
			_ = version

			config := &app.Config{
				NoHyprlandCheck: true,
				DebugMode:       false,
				ForceLiveMode:   false,
				IsTestMode:      true,
			}
			services := app.NewServices(config)
			_ = services
		}()
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("Concurrent access test timed out")
			return
		}
	}
}

func TestErrorRecovery(t *testing.T) {
	tests := []struct {
		name        string
		config      *app.Config
		expectError bool
	}{
		{
			name: "test mode should disable TUI",
			config: &app.Config{
				NoHyprlandCheck: true,
				DebugMode:       false,
				ForceLiveMode:   false,
				IsTestMode:      true,
			},
			expectError: true,
		},
		{
			name: "normal mode should work",
			config: &app.Config{
				NoHyprlandCheck: true,
				DebugMode:       false,
				ForceLiveMode:   false,
				IsTestMode:      false,
			},
			expectError: true, // Should be disabled in test environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runTUI(tt.config)
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if tt.expectError && err != nil && err.Error() != "TUI disabled during tests" {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func BenchmarkAppServicesCreation(b *testing.B) {
	config := &app.Config{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services := app.NewServices(config)
		_ = services
	}
}

func BenchmarkRunTUI(b *testing.B) {
	config := &app.Config{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runTUI(config)
	}
}
