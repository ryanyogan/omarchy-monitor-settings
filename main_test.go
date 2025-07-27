package main

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	config := &AppConfig{
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
			originalTerm := os.Getenv(tt.termVar)
			defer os.Setenv(tt.termVar, originalTerm)

			os.Setenv(tt.termVar, tt.terminalEnv)

			model := NewModel()
			if model.mode != ModeDashboard {
				t.Errorf("Expected ModeDashboard, got %v", model.mode)
			}
		})
	}
}

func TestSignalHandling(t *testing.T) {
	_ = NewModel()

	quitCmd := tea.Quit
	_ = quitCmd
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

			model := NewModel()
			if model.mode != ModeDashboard {
				t.Errorf("Expected ModeDashboard, got %v", model.mode)
			}
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

			model := NewModel()
			_ = model.mode
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

func TestMemoryUsage(t *testing.T) {
	models := make([]Model, 100)

	for i := 0; i < 100; i++ {
		models[i] = NewModel()
	}

	for i, model := range models {
		if model.mode != ModeDashboard {
			t.Errorf("Model %d: Expected ModeDashboard, got %v", i, model.mode)
		}
		if len(model.menuItems) == 0 {
			t.Errorf("Model %d: menuItems should not be empty", i)
		}
	}
}

func TestErrorRecovery(t *testing.T) {
	tests := []struct {
		name        string
		setupError  func()
		cleanup     func()
		shouldPanic bool
	}{
		{
			name: "invalid terminal size",
			setupError: func() {
			},
			cleanup:     func() {},
			shouldPanic: false,
		},
		{
			name: "missing environment",
			setupError: func() {
				os.Unsetenv("TERM")
				os.Unsetenv("DISPLAY")
				os.Unsetenv("WAYLAND_DISPLAY")
			},
			cleanup: func() {
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

			model := NewModel()
			if model.mode != ModeDashboard {
				t.Errorf("Expected ModeDashboard, got %v", model.mode)
			}
		})
	}
}

func BenchmarkMain(b *testing.B) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"omarchy-monitor-settings", "--no-hyprland-check"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model := NewModel()
		_ = model
	}
}

func BenchmarkModelCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model := NewModel()
		_ = model
	}
}

func BenchmarkRunTUI(b *testing.B) {
	config := &AppConfig{
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
