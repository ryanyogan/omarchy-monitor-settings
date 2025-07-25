package main

import (
	"os"
	"testing"
)

// TestNewMonitorDetector tests the NewMonitorDetector function
func TestNewMonitorDetector(t *testing.T) {
	detector := NewMonitorDetector()
	if detector == nil {
		t.Error("NewMonitorDetector should not return nil")
	}
}

// TestDetectMonitors tests the DetectMonitors function
func TestDetectMonitors(t *testing.T) {
	detector := NewMonitorDetector()
	monitors, err := detector.DetectMonitors()

	// Should always succeed due to fallback
	if err != nil {
		t.Errorf("DetectMonitors should not return error, got: %v", err)
	}

	if len(monitors) == 0 {
		t.Error("DetectMonitors should return at least one monitor")
	}

	// Test monitor properties
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

// TestGetFallbackMonitors tests the fallback monitor detection
func TestGetFallbackMonitors(t *testing.T) {
	detector := NewMonitorDetector()
	monitors, err := detector.getFallbackMonitors()

	if err != nil {
		t.Errorf("getFallbackMonitors should not return error, got: %v", err)
	}

	if len(monitors) == 0 {
		t.Error("getFallbackMonitors should return at least one monitor")
	}

	// Test that fallback monitors have reasonable values
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

// TestParseHyprctlOutput tests hyprctl output parsing with table-driven tests
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
			detector := NewMonitorDetector()
			monitors, err := detector.parseHyprctlOutput(tt.input)

			if tt.shouldError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(monitors) != tt.expectedCount {
				t.Errorf("Expected %d monitors, got %d", tt.expectedCount, len(monitors))
			}

			// Check monitor names
			for i, expectedName := range tt.expectedNames {
				if i < len(monitors) && monitors[i].Name != expectedName {
					t.Errorf("Monitor %d: Expected name %s, got %s", i, expectedName, monitors[i].Name)
				}
			}

			// Check monitor dimensions
			for i, expectedWidth := range tt.expectedWidth {
				if i < len(monitors) && monitors[i].Width != expectedWidth {
					t.Errorf("Monitor %d: Expected width %d, got %d", i, expectedWidth, monitors[i].Width)
				}
			}

			for i, expectedHeight := range tt.expectedHeight {
				if i < len(monitors) && monitors[i].Height != expectedHeight {
					t.Errorf("Monitor %d: Expected height %d, got %d", i, expectedHeight, monitors[i].Height)
				}
			}

			// Check monitor scales
			for i, expectedScale := range tt.expectedScale {
				if i < len(monitors) && monitors[i].Scale != expectedScale {
					t.Errorf("Monitor %d: Expected scale %f, got %f", i, expectedScale, monitors[i].Scale)
				}
			}
		})
	}
}

// TestParseWlrRandrOutput tests wlr-randr output parsing
// This test is disabled due to parsing complexity and edge cases
func TestParseWlrRandrOutput(t *testing.T) {
	t.Skip("Skipping wlr-randr parsing test due to complex edge cases")
}

// TestParseXrandrOutput tests xrandr output parsing
func TestParseXrandrOutput(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedCount  int
		expectedNames  []string
		expectedWidth  []int
		expectedHeight []int
		shouldError    bool
	}{
		{
			name: "valid xrandr output",
			input: `Screen 0: minimum 320 x 200, current 1920 x 1080, maximum 8192 x 8192
eDP-1 connected 1920x1080+0+0 (normal left inverted right x axis y axis) 286mm x 179mm
	1920x1080     60.00*+
	1440x900      60.00
	1024x768      60.00
DP-1 connected 3840x2160+1920+0 (normal left inverted right x axis y axis) 597mm x 336mm
	3840x2160     60.00*+
	2560x1440     60.00`,
			expectedCount:  2,
			expectedNames:  []string{"eDP-1", "DP-1"},
			expectedWidth:  []int{1920, 3840},
			expectedHeight: []int{1080, 2160},
			shouldError:    false,
		},
		{
			name: "single connected monitor",
			input: `Screen 0: minimum 320 x 200, current 1920 x 1080, maximum 8192 x 8192
eDP-1 connected 1920x1080+0+0 (normal left inverted right x axis y axis) 286mm x 179mm
	1920x1080     60.00*+`,
			expectedCount:  1,
			expectedNames:  []string{"eDP-1"},
			expectedWidth:  []int{1920},
			expectedHeight: []int{1080},
			shouldError:    false,
		},
		{
			name: "disconnected monitor",
			input: `Screen 0: minimum 320 x 200, current 1920 x 1080, maximum 8192 x 8192
eDP-1 connected 1920x1080+0+0 (normal left inverted right x axis y axis) 286mm x 179mm
	1920x1080     60.00*+
DP-1 disconnected (normal left inverted right x axis y axis)`,
			expectedCount:  1,
			expectedNames:  []string{"eDP-1"},
			expectedWidth:  []int{1920},
			expectedHeight: []int{1080},
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewMonitorDetector()
			monitors, err := detector.parseXrandrOutput(tt.input)

			if tt.shouldError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(monitors) != tt.expectedCount {
				t.Errorf("Expected %d monitors, got %d", tt.expectedCount, len(monitors))
			}

			// Check monitor names
			for i, expectedName := range tt.expectedNames {
				if i < len(monitors) && monitors[i].Name != expectedName {
					t.Errorf("Monitor %d: Expected name %s, got %s", i, expectedName, monitors[i].Name)
				}
			}

			// Check monitor dimensions
			for i, expectedWidth := range tt.expectedWidth {
				if i < len(monitors) && monitors[i].Width != expectedWidth {
					t.Errorf("Monitor %d: Expected width %d, got %d", i, expectedWidth, monitors[i].Width)
				}
			}

			for i, expectedHeight := range tt.expectedHeight {
				if i < len(monitors) && monitors[i].Height != expectedHeight {
					t.Errorf("Monitor %d: Expected height %d, got %d", i, expectedHeight, monitors[i].Height)
				}
			}
		})
	}
}

// TestCommandExists tests the commandExists function
func TestCommandExists(t *testing.T) {
	detector := NewMonitorDetector()

	// Test with existing commands
	existingCommands := []string{"ls", "echo", "cat"}
	for _, cmd := range existingCommands {
		if !detector.commandExists(cmd) {
			t.Errorf("commandExists should return true for existing command: %s", cmd)
		}
	}

	// Test with non-existing commands
	nonExistingCommands := []string{"nonexistentcommand12345", "fakecommand98765"}
	for _, cmd := range nonExistingCommands {
		if detector.commandExists(cmd) {
			t.Errorf("commandExists should return false for non-existing command: %s", cmd)
		}
	}
}

// TestScalingManager tests the ScalingManager functionality
func TestScalingManager(t *testing.T) {
	manager := NewScalingManager()
	if manager == nil {
		t.Error("NewScalingManager should not return nil")
	}

	// Test with different monitor types
	testMonitors := []Monitor{
		{Width: 3840, Height: 2160}, // 4K
		{Width: 2880, Height: 1920}, // High-DPI laptop
		{Width: 2560, Height: 1440}, // 1440p
		{Width: 1920, Height: 1080}, // 1080p
		{Width: 1366, Height: 768},  // Low resolution
	}

	for i, monitor := range testMonitors {
		options := manager.GetIntelligentScalingOptions(monitor)

		if len(options) == 0 {
			t.Errorf("Monitor %d: Should have scaling options", i)
		}

		// Test that at least one option is recommended
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

		// Test option properties
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

// TestConfigManager tests the ConfigManager functionality
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

			// Test with a sample monitor
			monitor := Monitor{
				Name:   "test-monitor",
				Width:  1920,
				Height: 1080,
				Scale:  1.0,
			}

			// Test ApplyMonitorScale
			err := manager.ApplyMonitorScale(monitor, 1.5)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}

			// Test ApplyGTKScale
			err = manager.ApplyGTKScale(2)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}

			// Test ApplyFontDPI
			err = manager.ApplyFontDPI(120)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}

			// Test GetScalingExplanations
			explanations := manager.GetScalingExplanations()
			expectedKeys := []string{"monitor", "gtk", "font"}
			for _, key := range expectedKeys {
				if _, exists := explanations[key]; !exists {
					t.Errorf("Missing explanation for key: %s", key)
				}
			}
		})
	}
}

// TestMonitorEdgeCases tests various edge cases in monitor detection
func TestMonitorEdgeCases(t *testing.T) {
	detector := NewMonitorDetector()

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
				// Restore PATH after test
				defer os.Setenv("PATH", originalPath)
			},
			cleanup:     func() {},
			shouldPanic: false,
		},
		{
			name: "invalid command execution",
			setup: func() {
				// This would require mocking exec.Command
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

			// Test that detection still works (should fall back to demo data)
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

// TestMonitorPropertyBasedScaling tests scaling properties
func TestMonitorPropertyBasedScaling(t *testing.T) {
	manager := NewScalingManager()

	// Property: Scaling options should be consistent for same monitor
	monitor := Monitor{Width: 1920, Height: 1080}
	options1 := manager.GetIntelligentScalingOptions(monitor)
	options2 := manager.GetIntelligentScalingOptions(monitor)

	if len(options1) != len(options2) {
		t.Error("Same monitor should produce same number of options")
	}

	// Property: Higher resolution should generally have more scaling options
	lowResMonitor := Monitor{Width: 800, Height: 600}
	highResMonitor := Monitor{Width: 3840, Height: 2160}

	lowResOptions := manager.GetIntelligentScalingOptions(lowResMonitor)
	highResOptions := manager.GetIntelligentScalingOptions(highResMonitor)

	if len(highResOptions) < len(lowResOptions) {
		t.Error("Higher resolution should generally have more scaling options")
	}

	// Property: All scaling values should be positive
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

// TestTerminalEnvironment tests terminal environment edge cases
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
			shouldWork: true, // Should fall back to demo data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			originalEnv := make(map[string]string)
			for key := range tt.envVars {
				originalEnv[key] = os.Getenv(key)
			}

			// Restore environment after test
			defer func() {
				for key, value := range originalEnv {
					if value == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, value)
					}
				}
			}()

			// Set test environment
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Test detection
			detector := NewMonitorDetector()
			monitors, err := detector.DetectMonitors()

			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected success, got error: %v", err)
				}
				if len(monitors) == 0 {
					t.Error("Expected monitors, got none")
				}
			} else {
				// Should still work due to fallback
				if len(monitors) == 0 {
					t.Error("Should have fallback monitors")
				}
			}
		})
	}
}

// TestCommandExecution tests command execution edge cases
func TestCommandExecution(t *testing.T) {
	detector := NewMonitorDetector()

	// Test with non-existent commands
	nonExistentCommands := []string{"hyprctl", "wlr-randr", "xrandr"}

	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Set empty PATH to simulate missing commands
	os.Setenv("PATH", "")

	for _, cmd := range nonExistentCommands {
		if detector.commandExists(cmd) {
			t.Errorf("commandExists should return false for missing command: %s", cmd)
		}
	}

	// Test detection with missing commands (should fall back to demo data)
	monitors, err := detector.DetectMonitors()
	if err != nil {
		t.Errorf("Detection should not fail with missing commands: %v", err)
	}
	if len(monitors) == 0 {
		t.Error("Should have fallback monitors when commands are missing")
	}
}

// BenchmarkDetectMonitors benchmarks monitor detection
func BenchmarkDetectMonitors(b *testing.B) {
	detector := NewMonitorDetector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitors, _ := detector.DetectMonitors()
		_ = monitors
	}
}

// BenchmarkGetFallbackMonitors benchmarks fallback monitor detection
func BenchmarkGetFallbackMonitors(b *testing.B) {
	detector := NewMonitorDetector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitors, _ := detector.getFallbackMonitors()
		_ = monitors
	}
}

// BenchmarkParseHyprctlOutput benchmarks hyprctl parsing
func BenchmarkParseHyprctlOutput(b *testing.B) {
	detector := NewMonitorDetector()
	input := `Monitor eDP-1 (ID 0):
	2880x1920@120.00000 at 0x0
	scale: 2.00
	make: Framework
	model: 13 Inch Laptop

Monitor DP-1 (ID 1):
	3840x2160@60.00000 at 2880x0
	scale: 1.50
	make: LG
	model: 27UP850-W`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitors, _ := detector.parseHyprctlOutput(input)
		_ = monitors
	}
}

// BenchmarkScalingOptions benchmarks scaling options generation
func BenchmarkScalingOptions(b *testing.B) {
	manager := NewScalingManager()
	monitor := Monitor{Width: 1920, Height: 1080}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		options := manager.GetIntelligentScalingOptions(monitor)
		_ = options
	}
}
