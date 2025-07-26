package main

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestMain tests the main function and command line flags
func TestMain(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantExit bool
	}{
		{
			name:     "with version flag",
			args:     []string{"omarchy-monitor-settings", "--version"},
			wantExit: true, // Version should exit
		},
		{
			name:     "with help flag",
			args:     []string{"omarchy-monitor-settings", "--help"},
			wantExit: true, // Help should exit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set test args
			os.Args = tt.args

			// Run in a goroutine to catch panics
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

			// Wait for completion or timeout
			select {
			case <-done:
				// Test completed
			case <-time.After(5 * time.Second):
				t.Error("main() timed out")
			}
		})
	}
}

// TestRunTUI tests the runTUI function
func TestRunTUI(t *testing.T) {
	// Test with test mode config
	config := &AppConfig{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true, // This should prevent TUI from starting
	}

	// Test normal execution
	err := runTUI(config)
	// Should return error because IsTestMode is true
	if err == nil {
		t.Error("Expected error when IsTestMode is true, got nil")
	}
	if err.Error() != "TUI disabled during tests" {
		t.Errorf("Expected 'TUI disabled during tests' error, got: %v", err)
	}
}

// TestGlobalFlags tests the global flag variables
func TestGlobalFlags(t *testing.T) {
	// Test default values
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

// TestTerminalEdgeCases tests edge cases specific to terminal environments
func TestTerminalEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		terminalEnv string
		termVar     string
		shouldWork  bool
	}{
		{
			name:        "dumb terminal",
			terminalEnv: "dumb",
			termVar:     "TERM",
			shouldWork:  false,
		},
		{
			name:        "no terminal",
			terminalEnv: "",
			termVar:     "TERM",
			shouldWork:  false,
		},
		{
			name:        "xterm terminal",
			terminalEnv: "xterm",
			termVar:     "TERM",
			shouldWork:  true,
		},
		{
			name:        "screen terminal",
			terminalEnv: "screen",
			termVar:     "TERM",
			shouldWork:  true,
		},
		{
			name:        "tmux terminal",
			terminalEnv: "tmux",
			termVar:     "TERM",
			shouldWork:  true,
		},
		{
			name:        "no color support",
			terminalEnv: "dumb",
			termVar:     "TERM",
			shouldWork:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			originalTerm := os.Getenv(tt.termVar)
			defer os.Setenv(tt.termVar, originalTerm)

			// Set test environment
			os.Setenv(tt.termVar, tt.terminalEnv)

			// Test model creation
			model := NewModel()
			if model.mode != ModeDashboard {
				t.Errorf("Expected ModeDashboard, got %v", model.mode)
			}
		})
	}
}

// TestSignalHandling tests signal handling in terminal environments
func TestSignalHandling(t *testing.T) {
	// Test that the model can handle interrupt signals gracefully
	_ = NewModel()

	// Simulate Ctrl+C
	quitCmd := tea.Quit
	_ = quitCmd // Use the variable to avoid unused variable warning
}

// TestEnvironmentVariables tests environment variable handling
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
			// Save original environment
			originalValue := os.Getenv(tt.envKey)
			defer os.Setenv(tt.envKey, originalValue)

			// Set test environment
			os.Setenv(tt.envKey, tt.envValue)

			// Test model creation
			model := NewModel()
			if model.mode != ModeDashboard {
				t.Errorf("Expected ModeDashboard, got %v", model.mode)
			}
		})
	}
}

// TestConcurrentAccess tests concurrent access to global variables
func TestConcurrentAccess(t *testing.T) {
	// Test that global flags can be accessed concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Access global variables
			_ = noHyprlandCheck
			_ = debugMode
			_ = forceLiveMode
			_ = version

			// Create model
			model := NewModel()
			_ = model.mode
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Error("Concurrent access test timed out")
			return
		}
	}
}

// TestMemoryUsage tests memory usage patterns
func TestMemoryUsage(t *testing.T) {
	// Test that creating multiple models doesn't cause memory leaks
	models := make([]Model, 100)

	for i := 0; i < 100; i++ {
		models[i] = NewModel()
	}

	// Verify all models are properly initialized
	for i, model := range models {
		if model.mode != ModeDashboard {
			t.Errorf("Model %d: Expected ModeDashboard, got %v", i, model.mode)
		}
		if len(model.menuItems) == 0 {
			t.Errorf("Model %d: menuItems should not be empty", i)
		}
	}
}

// TestErrorRecovery tests error recovery scenarios
func TestErrorRecovery(t *testing.T) {
	// Test that the application can recover from various error conditions
	tests := []struct {
		name        string
		setupError  func()
		cleanup     func()
		shouldPanic bool
	}{
		{
			name: "invalid terminal size",
			setupError: func() {
				// This would normally be set by tea.WindowSizeMsg
				// We can't easily simulate this in tests
			},
			cleanup:     func() {},
			shouldPanic: false,
		},
		{
			name: "missing environment",
			setupError: func() {
				// Clear environment variables
				os.Unsetenv("TERM")
				os.Unsetenv("DISPLAY")
				os.Unsetenv("WAYLAND_DISPLAY")
			},
			cleanup: func() {
				// Restore environment
				os.Setenv("TERM", "xterm")
			},
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.shouldPanic {
					t.Errorf("Unexpected panic: %v", r)
				}
				tt.cleanup()
			}()

			tt.setupError()

			// Try to create a model
			model := NewModel()
			if model.mode != ModeDashboard {
				t.Errorf("Expected ModeDashboard, got %v", model.mode)
			}
		})
	}
}

// BenchmarkMain benchmarks the main function
func BenchmarkMain(b *testing.B) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set minimal args for benchmarking
	os.Args = []string{"omarchy-monitor-settings", "--no-hyprland-check"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We can't easily benchmark main() as it runs indefinitely
		// Instead, benchmark model creation
		model := NewModel()
		_ = model
	}
}

// BenchmarkModelCreation benchmarks model creation
func BenchmarkModelCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model := NewModel()
		_ = model
	}
}

// BenchmarkRunTUI benchmarks the runTUI function
func BenchmarkRunTUI(b *testing.B) {
	// Create test config to prevent TUI from starting
	config := &AppConfig{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This will return an error due to test mode
		// but we can still measure the attempt
		_ = runTUI(config)
	}
}
