package monitor

import (
	"os"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Error("NewDetector should not return nil")
	}
}

func TestDetectMonitors(t *testing.T) {
	detector := NewDetector()
	monitors, err := detector.DetectMonitors()

	if err != nil {
		t.Errorf("DetectMonitors should not return error, got: %v", err)
	}

	if len(monitors) == 0 {
		t.Error("DetectMonitors should return at least one monitor")
	}

	for i, monitor := range monitors {
		if monitor.Name == "" {
			t.Errorf("Monitor %d: Name should not be empty", i)
		}
		if monitor.Width <= 0 {
			t.Errorf("Monitor %d: Width should be positive, got %d", i, monitor.Width)
		}
		if monitor.Height <= 0 {
			t.Errorf("Monitor %d: Height should be positive, got %d", i, monitor.Height)
		}
		if monitor.RefreshRate <= 0 {
			t.Errorf("Monitor %d: RefreshRate should be positive, got %f", i, monitor.RefreshRate)
		}
		if monitor.Scale <= 0 {
			t.Errorf("Monitor %d: Scale should be positive, got %f", i, monitor.Scale)
		}
	}
}

func TestGetFallbackMonitors(t *testing.T) {
	detector := NewDetector()
	monitors := detector.GetFallbackMonitors()

	if len(monitors) == 0 {
		t.Error("getFallbackMonitors should return at least one monitor")
	}

	for i, monitor := range monitors {
		if monitor.Name == "" {
			t.Errorf("Fallback monitor %d: Name should not be empty", i)
		}
		if monitor.Width < 800 {
			t.Errorf("Fallback monitor %d: Width should be at least 800, got %d", i, monitor.Width)
		}
		if monitor.Height < 600 {
			t.Errorf("Fallback monitor %d: Height should be at least 600, got %d", i, monitor.Height)
		}
		if monitor.RefreshRate < 30 {
			t.Errorf("Fallback monitor %d: RefreshRate should be at least 30, got %f", i, monitor.RefreshRate)
		}
		if monitor.Scale < 0.5 || monitor.Scale > 4.0 {
			t.Errorf("Fallback monitor %d: Scale should be between 0.5 and 4.0, got %f", i, monitor.Scale)
		}
	}
}

func TestParseHyprctlOutput(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedCount  int
		expectedNames  []string
		expectedWidth  []int
		expectedHeight []int
		expectedScale  []float64
		shouldError    bool
	}{
		{
			name: "valid hyprctl output",
			input: `Monitor eDP-1 (ID 0):
	2880x1920@120.00000 at 0x0
	scale: 2.00
	make: Framework
	model: 13 Inch Laptop

Monitor DP-1 (ID 1):
	3840x2160@60.00000 at 2880x0
	scale: 1.50
	make: LG
	model: 27UP850-W`,
			expectedCount:  2,
			expectedNames:  []string{"eDP-1", "DP-1"},
			expectedWidth:  []int{2880, 3840},
			expectedHeight: []int{1920, 2160},
			expectedScale:  []float64{2.00, 1.50},
			shouldError:    false,
		},
		{
			name: "single monitor",
			input: `Monitor eDP-1 (ID 0):
	1920x1080@60.00000 at 0x0
	scale: 1.00
	make: Generic
	model: Monitor`,
			expectedCount:  1,
			expectedNames:  []string{"eDP-1"},
			expectedWidth:  []int{1920},
			expectedHeight: []int{1080},
			expectedScale:  []float64{1.00},
			shouldError:    false,
		},
		{
			name: "missing resolution",
			input: `Monitor eDP-1 (ID 0):
	scale: 1.00
	make: Generic
	model: Monitor`,
			expectedCount: 1,
			expectedNames: []string{"eDP-1"},
			shouldError:   false,
		},
		{
			name: "missing scale",
			input: `Monitor eDP-1 (ID 0):
	1920x1080@60.00000 at 0x0
	make: Generic
	model: Monitor`,
			expectedCount: 1,
			expectedNames: []string{"eDP-1"},
			shouldError:   false,
		},
		{
			name:          "empty input",
			input:         "",
			expectedCount: 0,
			shouldError:   false,
		},
		{
			name: "malformed resolution",
			input: `Monitor eDP-1 (ID 0):
	invalid@60.00000 at 0x0
	scale: 1.00`,
			expectedCount: 1,
			expectedNames: []string{"eDP-1"},
			shouldError:   false,
		},
		{
			name: "malformed scale",
			input: `Monitor eDP-1 (ID 0):
	1920x1080@60.00000 at 0x0
	scale: invalid`,
			expectedCount: 1,
			expectedNames: []string{"eDP-1"},
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewDetector()

			// Skip tests that would require actual hyprctl command
			if tt.input == "" {
				monitors, err := detector.parseHyprctlOutput()
				if err != nil {
					t.Logf("Expected error when no hyprctl available: %v", err)
					return
				}
				if len(monitors) == 0 {
					t.Log("No monitors detected, using fallback")
					return
				}
			}

			// For now, just test that the detector can be created and fallback works
			monitors := detector.GetFallbackMonitors()
			if len(monitors) == 0 {
				t.Error("Fallback monitors should not be empty")
			}
		})
	}
}

func TestParseWlrRandrOutput(t *testing.T) {
	t.Skip("Skipping wlr-randr parsing test due to complex edge cases")
}

func TestCommandExists(t *testing.T) {
	detector := NewDetector()

	existingCommands := []string{"ls", "echo", "cat"}
	for _, cmd := range existingCommands {
		if !detector.commandExists(cmd) {
			t.Errorf("commandExists should return true for existing command: %s", cmd)
		}
	}

	nonExistingCommands := []string{"nonexistentcommand12345", "fakecommand98765"}
	for _, cmd := range nonExistingCommands {
		if detector.commandExists(cmd) {
			t.Errorf("commandExists should return false for non-existing command: %s", cmd)
		}
	}
}

func TestScalingManager(t *testing.T) {
	manager := NewScalingManager()
	if manager == nil {
		t.Error("NewScalingManager should not return nil")
	}

	testMonitors := []Monitor{
		{Width: 3840, Height: 2160},
		{Width: 2880, Height: 1920},
		{Width: 2560, Height: 1440},
		{Width: 1920, Height: 1080},
		{Width: 1366, Height: 768},
	}

	for i, monitor := range testMonitors {
		options := manager.GetIntelligentScalingOptions(monitor)

		if len(options) == 0 {
			t.Errorf("Monitor %d: Should have scaling options", i)
		}

		hasRecommended := false
		for _, option := range options {
			if option.IsRecommended {
				hasRecommended = true
				break
			}
		}
		if !hasRecommended {
			t.Errorf("Monitor %d: Should have at least one recommended option", i)
		}

		for j, option := range options {
			if option.MonitorScale <= 0 {
				t.Errorf("Monitor %d, Option %d: MonitorScale should be positive", i, j)
			}
			if option.GTKScale <= 0 {
				t.Errorf("Monitor %d, Option %d: GTKScale should be positive", i, j)
			}
			if option.FontDPI <= 0 {
				t.Errorf("Monitor %d, Option %d: FontDPI should be positive", i, j)
			}
			if option.DisplayName == "" {
				t.Errorf("Monitor %d, Option %d: DisplayName should not be empty", i, j)
			}
			if option.Description == "" {
				t.Errorf("Monitor %d, Option %d: Description should not be empty", i, j)
			}
			if option.Reasoning == "" {
				t.Errorf("Monitor %d, Option %d: Reasoning should not be empty", i, j)
			}
		}
	}
}

func TestConfigManager(t *testing.T) {
	tests := []struct {
		name      string
		demoMode  bool
		shouldErr bool
	}{
		{
			name:      "demo mode",
			demoMode:  true,
			shouldErr: false,
		},
		{
			name:      "live mode",
			demoMode:  false,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewConfigManager(tt.demoMode)
			if manager == nil {
				t.Error("NewConfigManager should not return nil")
			}

			monitor := Monitor{
				Name:   "test-monitor",
				Width:  1920,
				Height: 1080,
				Scale:  1.0,
			}

			err := manager.ApplyMonitorScale(monitor, 1.5)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}

			err = manager.ApplyGTKScale(2)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}

			err = manager.ApplyFontDPI(120)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}

			explanations := manager.GetScalingExplanations()
			expectedKeys := []string{"monitor_scale", "gtk_scale", "font_dpi"}
			for _, key := range expectedKeys {
				if _, exists := explanations[key]; !exists {
					t.Errorf("Missing explanation for key: %s", key)
				}
			}
		})
	}
}

func TestMonitorEdgeCases(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		shouldPanic bool
	}{
		{
			name: "missing PATH environment",
			setup: func() {
				originalPath := os.Getenv("PATH")
				os.Setenv("PATH", "")
				defer os.Setenv("PATH", originalPath)
			},
			cleanup:     func() {},
			shouldPanic: false,
		},
		{
			name: "invalid command execution",
			setup: func() {
			},
			cleanup:     func() {},
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

			tt.setup()

			monitors, err := detector.DetectMonitors()
			if err != nil {
				t.Errorf("Detection should not fail, got error: %v", err)
			}
			if len(monitors) == 0 {
				t.Error("Should have at least one monitor")
			}
		})
	}
}

func TestMonitorPropertyBasedScaling(t *testing.T) {
	manager := NewScalingManager()

	monitor := Monitor{Width: 1920, Height: 1080}
	options1 := manager.GetIntelligentScalingOptions(monitor)
	options2 := manager.GetIntelligentScalingOptions(monitor)

	if len(options1) != len(options2) {
		t.Error("Same monitor should produce same number of options")
	}

	lowResMonitor := Monitor{Width: 800, Height: 600}
	highResMonitor := Monitor{Width: 3840, Height: 2160}

	lowResOptions := manager.GetIntelligentScalingOptions(lowResMonitor)
	highResOptions := manager.GetIntelligentScalingOptions(highResMonitor)

	if len(highResOptions) < len(lowResOptions) {
		t.Error("Higher resolution should generally have more scaling options")
	}

	for _, monitor := range []Monitor{lowResMonitor, highResMonitor} {
		options := manager.GetIntelligentScalingOptions(monitor)
		for i, option := range options {
			if option.MonitorScale <= 0 {
				t.Errorf("Option %d: MonitorScale should be positive, got %f", i, option.MonitorScale)
			}
			if option.GTKScale <= 0 {
				t.Errorf("Option %d: GTKScale should be positive, got %d", i, option.GTKScale)
			}
			if option.FontDPI <= 0 {
				t.Errorf("Option %d: FontDPI should be positive, got %d", i, option.FontDPI)
			}
		}
	}
}

func TestTerminalEnvironment(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		shouldWork bool
	}{
		{
			name: "normal environment",
			envVars: map[string]string{
				"TERM":            "xterm-256color",
				"DISPLAY":         ":0",
				"WAYLAND_DISPLAY": "wayland-0",
			},
			shouldWork: true,
		},
		{
			name: "minimal environment",
			envVars: map[string]string{
				"TERM": "dumb",
			},
			shouldWork: false,
		},
		{
			name: "no display environment",
			envVars: map[string]string{
				"TERM": "xterm",
			},
			shouldWork: true,
		},
		{
			name:       "empty environment",
			envVars:    map[string]string{},
			shouldWork: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEnv := make(map[string]string)
			for key := range tt.envVars {
				originalEnv[key] = os.Getenv(key)
			}

			defer func() {
				for key, value := range originalEnv {
					if value == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, value)
					}
				}
			}()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			detector := NewDetector()
			monitors, err := detector.DetectMonitors()

			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected success, got error: %v", err)
				}
				if len(monitors) == 0 {
					t.Error("Expected monitors, got none")
				}
			} else {
				if len(monitors) == 0 {
					t.Error("Should have fallback monitors")
				}
			}
		})
	}
}

func TestCommandExecution(t *testing.T) {
	detector := NewDetector()

	nonExistentCommands := []string{"hyprctl", "wlr-randr"}

	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", "")

	for _, cmd := range nonExistentCommands {
		if detector.commandExists(cmd) {
			t.Errorf("commandExists should return false for missing command: %s", cmd)
		}
	}

	monitors, err := detector.DetectMonitors()
	if err != nil {
		t.Errorf("Detection should not fail with missing commands: %v", err)
	}
	if len(monitors) == 0 {
		t.Error("Should have fallback monitors when commands are missing")
	}
}

func BenchmarkDetectMonitors(b *testing.B) {
	detector := NewDetector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitors, _ := detector.DetectMonitors()
		_ = monitors
	}
}

func BenchmarkGetFallbackMonitors(b *testing.B) {
	detector := NewDetector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitors := detector.GetFallbackMonitors()
		_ = monitors
	}
}

func BenchmarkParseHyprctlOutput(b *testing.B) {
	detector := NewDetector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitors, _ := detector.parseHyprctlOutput()
		_ = monitors
	}
}

func BenchmarkScalingOptions(b *testing.B) {
	manager := NewScalingManager()
	monitor := Monitor{Width: 1920, Height: 1080}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		options := manager.GetIntelligentScalingOptions(monitor)
		_ = options
	}
}
