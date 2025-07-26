package main

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestNewModel tests the NewModel function
func TestNewModel(t *testing.T) {
	model := NewModel()

	// Test basic initialization
	if model.mode != ModeDashboard {
		t.Errorf("Expected ModeDashboard, got %v", model.mode)
	}
	if model.selectedOption != 0 {
		t.Errorf("Expected selectedOption 0, got %d", model.selectedOption)
	}
	if len(model.menuItems) == 0 {
		t.Error("menuItems should not be empty")
	}
	if model.manualMonitorScale != 1.0 {
		t.Errorf("Expected manualMonitorScale 1.0, got %f", model.manualMonitorScale)
	}
	if model.manualGTKScale != 1 {
		t.Errorf("Expected manualGTKScale 1, got %d", model.manualGTKScale)
	}
	if model.manualFontDPI != 96 {
		t.Errorf("Expected manualFontDPI 96, got %d", model.manualFontDPI)
	}
	if model.selectedManualControl != 0 {
		t.Errorf("Expected selectedManualControl 0, got %d", model.selectedManualControl)
	}
}

// TestModelInit tests the Init method
func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

// TestModelUpdate tests the Update method with various messages
func TestModelUpdate(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		expected AppMode
	}{
		{
			name:     "window size message",
			msg:      tea.WindowSizeMsg{Width: 80, Height: 24},
			expected: ModeDashboard,
		},
		{
			name:     "key message quit",
			msg:      tea.KeyMsg{Type: tea.KeyCtrlC},
			expected: ModeDashboard,
		},
		{
			name:     "key message up",
			msg:      tea.KeyMsg{Type: tea.KeyUp},
			expected: ModeDashboard,
		},
		{
			name:     "key message down",
			msg:      tea.KeyMsg{Type: tea.KeyDown},
			expected: ModeDashboard,
		},
		{
			name:     "key message enter",
			msg:      tea.KeyMsg{Type: tea.KeyEnter},
			expected: ModeDashboard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			updatedModel, cmd := model.Update(tt.msg)

			// Type assert to access model fields
			if model, ok := updatedModel.(Model); ok {
				if model.mode != tt.expected {
					t.Errorf("Expected mode %v, got %v", tt.expected, model.mode)
				}
			}

			// Test that we can chain updates
			if cmd != nil {
				// This would normally execute the command
				_ = cmd
			}
		})
	}
}

// TestHandleKeyPress tests key press handling with table-driven tests
func TestHandleKeyPress(t *testing.T) {
	tests := []struct {
		name           string
		initialMode    AppMode
		initialOption  int
		key            string
		expectedMode   AppMode
		expectedOption int
		shouldQuit     bool
	}{
		{
			name:           "quit with q",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "q",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     true,
		},
		{
			name:           "quit with ctrl+c",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "ctrl+c",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     true,
		},
		{
			name:           "navigate up",
			initialMode:    ModeDashboard,
			initialOption:  1,
			key:            "up",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "navigate down",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "down",
			expectedMode:   ModeDashboard,
			expectedOption: 1,
			shouldQuit:     false,
		},
		{
			name:           "navigate with k",
			initialMode:    ModeDashboard,
			initialOption:  1,
			key:            "k",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "navigate with j",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "j",
			expectedMode:   ModeDashboard,
			expectedOption: 1,
			shouldQuit:     false,
		},
		{
			name:           "select dashboard",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "enter",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "select monitor selection",
			initialMode:    ModeDashboard,
			initialOption:  1,
			key:            "enter",
			expectedMode:   ModeMonitorSelection,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "select smart scaling",
			initialMode:    ModeDashboard,
			initialOption:  2,
			key:            "enter",
			expectedMode:   ModeScalingOptions,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "select manual scaling",
			initialMode:    ModeDashboard,
			initialOption:  3,
			key:            "enter",
			expectedMode:   ModeManualScaling,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "select settings",
			initialMode:    ModeDashboard,
			initialOption:  4,
			key:            "enter",
			expectedMode:   ModeSettings,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "select help",
			initialMode:    ModeDashboard,
			initialOption:  5,
			key:            "enter",
			expectedMode:   ModeHelp,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "select exit",
			initialMode:    ModeDashboard,
			initialOption:  6,
			key:            "enter",
			expectedMode:   ModeDashboard,
			expectedOption: 6,
			shouldQuit:     true,
		},
		{
			name:           "help key",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "h",
			expectedMode:   ModeHelp,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "help key with question mark",
			initialMode:    ModeDashboard,
			initialOption:  0,
			key:            "?",
			expectedMode:   ModeHelp,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "escape from help",
			initialMode:    ModeHelp,
			initialOption:  0,
			key:            "esc",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     false,
		},
		{
			name:           "escape from manual scaling",
			initialMode:    ModeManualScaling,
			initialOption:  0,
			key:            "esc",
			expectedMode:   ModeDashboard,
			expectedOption: 0,
			shouldQuit:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.mode = tt.initialMode
			model.selectedOption = tt.initialOption

			// Create key message
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "up":
				keyMsg.Type = tea.KeyUp
			case "down":
				keyMsg.Type = tea.KeyDown
			case "enter":
				keyMsg.Type = tea.KeyEnter
			case "esc":
				keyMsg.Type = tea.KeyEscape
			case "ctrl+c":
				keyMsg.Type = tea.KeyCtrlC
			case "k":
				keyMsg.Type = tea.KeyRunes
				keyMsg.Runes = []rune("k")
			case "j":
				keyMsg.Type = tea.KeyRunes
				keyMsg.Runes = []rune("j")
			case "q":
				keyMsg.Type = tea.KeyRunes
				keyMsg.Runes = []rune("q")
			case "h":
				keyMsg.Type = tea.KeyRunes
				keyMsg.Runes = []rune("h")
			case "?":
				keyMsg.Type = tea.KeyRunes
				keyMsg.Runes = []rune("?")
			}

			updatedModel, cmd := model.handleKeyPress(keyMsg)

			// Type assert to access model fields
			if model, ok := updatedModel.(Model); ok {
				if model.mode != tt.expectedMode {
					t.Errorf("Expected mode %v, got %v", tt.expectedMode, model.mode)
				}
			}

			if tt.shouldQuit {
				if cmd == nil {
					t.Error("Expected quit command, got nil")
				}
			}
		})
	}
}

// TestManualScalingControls tests manual scaling control interactions
func TestManualScalingControls(t *testing.T) {
	tests := []struct {
		name                 string
		initialControl       int
		initialMonitorScale  float64
		initialGTKScale      int
		initialFontDPI       int
		key                  string
		expectedControl      int
		expectedMonitorScale float64
		expectedGTKScale     int
		expectedFontDPI      int
	}{
		{
			name:                 "switch to monitor scale control",
			initialControl:       1,
			initialMonitorScale:  1.0,
			key:                  "up",
			expectedControl:      0,
			expectedMonitorScale: 1.0,
		},
		{
			name:             "switch to GTK scale control",
			initialControl:   0,
			initialGTKScale:  1,
			key:              "down",
			expectedControl:  1,
			expectedGTKScale: 1,
		},
		{
			name:            "switch to font DPI control",
			initialControl:  1,
			initialFontDPI:  96,
			key:             "down",
			expectedControl: 2,
			expectedFontDPI: 96,
		},
		{
			name:                 "increase monitor scale",
			initialControl:       0,
			initialMonitorScale:  1.0,
			key:                  "right",
			expectedControl:      0,
			expectedMonitorScale: 1.25,
		},
		{
			name:                 "decrease monitor scale",
			initialControl:       0,
			initialMonitorScale:  1.5,
			key:                  "left",
			expectedControl:      0,
			expectedMonitorScale: 1.33333,
		},
		{
			name:             "increase GTK scale",
			initialControl:   1,
			initialGTKScale:  1,
			key:              "right",
			expectedControl:  1,
			expectedGTKScale: 2,
		},
		{
			name:             "decrease GTK scale",
			initialControl:   1,
			initialGTKScale:  2,
			key:              "left",
			expectedControl:  1,
			expectedGTKScale: 1,
		},
		{
			name:            "increase font DPI",
			initialControl:  2,
			initialFontDPI:  96,
			key:             "right",
			expectedControl: 2,
			expectedFontDPI: 108,
		},
		{
			name:            "decrease font DPI",
			initialControl:  2,
			initialFontDPI:  108,
			key:             "left",
			expectedControl: 2,
			expectedFontDPI: 96,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.mode = ModeManualScaling
			model.selectedManualControl = tt.initialControl
			model.manualMonitorScale = tt.initialMonitorScale
			model.manualGTKScale = tt.initialGTKScale
			model.manualFontDPI = tt.initialFontDPI

			// Create key message
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "up":
				keyMsg.Type = tea.KeyUp
			case "down":
				keyMsg.Type = tea.KeyDown
			}

			updatedModel, _ := model.handleKeyPress(keyMsg)

			// Type assert to access model fields
			if model, ok := updatedModel.(Model); ok {
				if model.selectedManualControl != tt.expectedControl {
					t.Errorf("Expected control %d, got %d", tt.expectedControl, model.selectedManualControl)
				}
				if model.manualMonitorScale != tt.expectedMonitorScale {
					t.Errorf("Expected monitor scale %f, got %f", tt.expectedMonitorScale, model.manualMonitorScale)
				}
				if model.manualGTKScale != tt.expectedGTKScale {
					t.Errorf("Expected GTK scale %d, got %d", tt.expectedGTKScale, model.manualGTKScale)
				}
				if model.manualFontDPI != tt.expectedFontDPI {
					t.Errorf("Expected font DPI %d, got %d", tt.expectedFontDPI, model.manualFontDPI)
				}
			}
		})
	}
}

// TestModelView tests the View method with various terminal sizes
func TestModelView(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		ready    bool
		mode     AppMode
		expected string
	}{
		{
			name:     "terminal too small",
			width:    40,
			height:   10,
			ready:    true,
			mode:     ModeDashboard,
			expected: "Terminal too small",
		},
		{
			name:     "not ready",
			width:    80,
			height:   24,
			ready:    false,
			mode:     ModeDashboard,
			expected: "Initializing",
		},
		{
			name:     "dashboard mode",
			width:    80,
			height:   24,
			ready:    true,
			mode:     ModeDashboard,
			expected: "Display Settings",
		},
		{
			name:     "help mode",
			width:    80,
			height:   24,
			ready:    true,
			mode:     ModeHelp,
			expected: "Help & Controls",
		},
		{
			name:     "settings mode",
			width:    80,
			height:   24,
			ready:    true,
			mode:     ModeSettings,
			expected: "Application Settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.width = tt.width
			model.height = tt.height
			model.ready = tt.ready
			model.mode = tt.mode

			view := model.View()

			if !strings.Contains(view, tt.expected) {
				t.Errorf("View should contain '%s', got: %s", tt.expected, view)
			}
		})
	}
}

// TestScalingOptions tests scaling options functionality
func TestScalingOptions(t *testing.T) {
	model := NewModel()

	// Ensure we have monitors
	if len(model.monitors) == 0 {
		t.Skip("No monitors available for testing")
	}

	// Test scaling options generation
	scalingManager := NewScalingManager()
	options := scalingManager.GetIntelligentScalingOptions(model.monitors[0])

	if len(options) == 0 {
		t.Error("Should have scaling options")
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
		t.Error("Should have at least one recommended option")
	}
}

// TestMonitorSelection tests monitor selection functionality
func TestMonitorSelection(t *testing.T) {
	model := NewModel()

	// Ensure we have monitors
	if len(model.monitors) == 0 {
		t.Skip("No monitors available for testing")
	}

	// Test monitor selection navigation
	model.mode = ModeMonitorSelection
	model.selectedMonitor = 0

	// Navigate down
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.handleKeyPress(keyMsg)

	// Type assert to access model fields
	if model, ok := updatedModel.(Model); ok {
		if len(model.monitors) > 1 && model.selectedMonitor != 1 {
			t.Errorf("Expected selectedMonitor 1, got %d", model.selectedMonitor)
		}

		// Navigate up
		keyMsg = tea.KeyMsg{Type: tea.KeyUp}
		updatedModel2, _ := model.handleKeyPress(keyMsg)

		if model2, ok := updatedModel2.(Model); ok {
			if model2.selectedMonitor != 0 {
				t.Errorf("Expected selectedMonitor 0, got %d", model2.selectedMonitor)
			}
		}
	}
}

// TestConfirmationFlow tests the confirmation flow
func TestConfirmationFlow(t *testing.T) {
	model := NewModel()

	// Ensure we have monitors
	if len(model.monitors) == 0 {
		t.Skip("No monitors available for testing")
	}

	// Set up for smart scaling confirmation
	model.mode = ModeScalingOptions
	model.selectedMonitor = 0
	model.selectedScalingOpt = 0

	// Select a scaling option
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.handleKeyPress(keyMsg)

	// Type assert to access model fields
	if model, ok := updatedModel.(Model); ok {
		if model.mode != ModeConfirmation {
			t.Errorf("Expected ModeConfirmation, got %v", model.mode)
		}

		// Cancel confirmation
		keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
		updatedModel2, _ := model.handleKeyPress(keyMsg)

		if model2, ok := updatedModel2.(Model); ok {
			if model2.mode != ModeScalingOptions {
				t.Errorf("Expected ModeScalingOptions, got %v", model2.mode)
			}
		}
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupModel  func() Model
		key         string
		shouldPanic bool
	}{
		{
			name: "empty monitors list",
			setupModel: func() Model {
				model := NewModel()
				model.monitors = []Monitor{}
				return model
			},
			key:         "enter",
			shouldPanic: false,
		},
		{
			name: "invalid monitor selection",
			setupModel: func() Model {
				model := NewModel()
				model.selectedMonitor = 999 // Invalid index
				return model
			},
			key:         "enter",
			shouldPanic: false,
		},
		{
			name: "invalid scaling option selection",
			setupModel: func() Model {
				model := NewModel()
				model.mode = ModeScalingOptions
				model.selectedScalingOpt = 999 // Invalid index
				return model
			},
			key:         "enter",
			shouldPanic: false,
		},
		{
			name: "boundary monitor scale values",
			setupModel: func() Model {
				model := NewModel()
				model.mode = ModeManualScaling
				model.selectedManualControl = 0
				model.manualMonitorScale = 0.1 // Very small value
				return model
			},
			key:         "down",
			shouldPanic: false,
		},
		{
			name: "boundary GTK scale values",
			setupModel: func() Model {
				model := NewModel()
				model.mode = ModeManualScaling
				model.selectedManualControl = 1
				model.manualGTKScale = 1 // Minimum value
				return model
			},
			key:         "down",
			shouldPanic: false,
		},
		{
			name: "boundary font DPI values",
			setupModel: func() Model {
				model := NewModel()
				model.mode = ModeManualScaling
				model.selectedManualControl = 2
				model.manualFontDPI = 72 // Minimum value
				return model
			},
			key:         "down",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.shouldPanic {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			model := tt.setupModel()

			// Create key message
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "enter":
				keyMsg.Type = tea.KeyEnter
			case "down":
				keyMsg.Type = tea.KeyDown
			}

			// This should not panic
			updatedModel, _ := model.handleKeyPress(keyMsg)
			_ = updatedModel
		})
	}
}

// TestPropertyBasedScaling tests scaling properties
func TestPropertyBasedScaling(t *testing.T) {
	// Property: Scaling values should always be positive
	model := NewModel()
	model.mode = ModeManualScaling

	// Test monitor scale bounds
	if model.manualMonitorScale <= 0 {
		t.Error("Monitor scale should be positive")
	}

	// Test GTK scale bounds
	if model.manualGTKScale <= 0 {
		t.Error("GTK scale should be positive")
	}

	// Test font DPI bounds
	if model.manualFontDPI <= 0 {
		t.Error("Font DPI should be positive")
	}

	// Property: After increasing and decreasing, we should get back to original value
	originalScale := model.manualMonitorScale

	// Increase scale
	keyMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.handleKeyPress(keyMsg)

	// Decrease scale - need to type assert first
	if model, ok := updatedModel.(Model); ok {
		keyMsg = tea.KeyMsg{Type: tea.KeyDown}
		finalModel, _ := model.handleKeyPress(keyMsg)

		// Type assert to access model fields
		if finalModel, ok := finalModel.(Model); ok {
			if finalModel.manualMonitorScale != originalScale {
				t.Errorf("Scale should return to original value, got %f, expected %f",
					finalModel.manualMonitorScale, originalScale)
			}
		}
	}
}

// TestTerminalSizeConstraints tests terminal size constraints
func TestTerminalSizeConstraints(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		shouldWork bool
	}{
		{
			name:       "minimum size",
			width:      80,
			height:     20,
			shouldWork: true,
		},
		{
			name:       "below minimum width",
			width:      40,
			height:     20,
			shouldWork: false,
		},
		{
			name:       "below minimum height",
			width:      80,
			height:     10,
			shouldWork: false,
		},
		{
			name:       "very large terminal",
			width:      200,
			height:     100,
			shouldWork: true,
		},
		{
			name:       "zero dimensions",
			width:      0,
			height:     0,
			shouldWork: false,
		},
		{
			name:       "negative dimensions",
			width:      -10,
			height:     -10,
			shouldWork: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.width = tt.width
			model.height = tt.height
			model.ready = true

			view := model.View()

			if tt.shouldWork {
				if strings.Contains(view, "Terminal too small") {
					t.Error("Terminal should be large enough")
				}
			} else {
				if !strings.Contains(view, "Terminal too small") {
					t.Error("Terminal should be too small")
				}
			}
		})
	}
}

// BenchmarkModelUpdate benchmarks the Update method
func BenchmarkModelUpdate(b *testing.B) {
	model := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedModel, _ := model.Update(msg)
		_ = updatedModel
	}
}

// BenchmarkModelView benchmarks the View method
func BenchmarkModelView(b *testing.B) {
	model := NewModel()
	model.width = 80
	model.height = 24
	model.ready = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		view := model.View()
		_ = view
	}
}

// BenchmarkKeyPress benchmarks key press handling
func BenchmarkKeyPress(b *testing.B) {
	model := NewModel()
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedModel, _ := model.handleKeyPress(keyMsg)
		_ = updatedModel
	}
}

// TestTerminalColorInheritance tests that the application correctly inherits terminal colors
func TestTerminalColorInheritance(t *testing.T) {
	t.Run("ANSI_color_definitions", func(t *testing.T) {
		// Test that all color variables use ANSI color codes that adapt to terminal themes
		testCases := []struct {
			name     string
			color    lipgloss.Color
			expected string
		}{
			{"colorBackground", colorBackground, ""}, // Terminal default
			{"colorSurface", colorSurface, "0"},      // ANSI black
			{"colorFloat", colorFloat, "8"},          // ANSI bright black
			{"colorForeground", colorForeground, ""}, // Terminal default
			{"colorComment", colorComment, "8"},      // ANSI bright black
			{"colorSubtle", colorSubtle, "7"},        // ANSI white
			{"colorBlue", colorBlue, "4"},            // ANSI blue
			{"colorCyan", colorCyan, "6"},            // ANSI cyan
			{"colorGreen", colorGreen, "2"},          // ANSI green
			{"colorYellow", colorYellow, "3"},        // ANSI yellow
			{"colorRed", colorRed, "1"},              // ANSI red
			{"colorMagenta", colorMagenta, "5"},      // ANSI magenta
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				actual := string(tc.color)
				if actual != tc.expected {
					t.Errorf("Expected %s to be ANSI color '%s', got '%s'", tc.name, tc.expected, actual)
				}
			})
		}
	})

	t.Run("terminal_theme_detection", func(t *testing.T) {
		// Test the getTerminalThemeInfo function
		themeInfo := getTerminalThemeInfo()

		// Should contain "Terminal Adaptive" and additional info
		if !strings.Contains(themeInfo, "Terminal Adaptive") {
			t.Errorf("Expected theme info to contain 'Terminal Adaptive', got: %s", themeInfo)
		}

		// Should contain profile information
		profiles := []string{"TrueColor", "256 Color", "16 Color", "Basic"}
		hasProfile := false
		for _, profile := range profiles {
			if strings.Contains(themeInfo, profile) {
				hasProfile = true
				break
			}
		}
		if !hasProfile {
			t.Errorf("Expected theme info to contain a color profile, got: %s", themeInfo)
		}

		// Should contain theme type
		themes := []string{"Dark", "Light"}
		hasTheme := false
		for _, theme := range themes {
			if strings.Contains(themeInfo, theme) {
				hasTheme = true
				break
			}
		}
		if !hasTheme {
			t.Errorf("Expected theme info to contain theme type (Dark/Light), got: %s", themeInfo)
		}
	})

	t.Run("style_color_usage", func(t *testing.T) {
		// Create a model and initialize styles
		config := &AppConfig{
			IsTestMode: true,
		}
		services := NewAppServices(config)
		m := NewModelWithServices(services)
		m.width = 100
		m.height = 30
		m.ready = true

		// Test that styles are properly using terminal-adaptive colors
		// Call View() to ensure width/height are used
		_ = m.View()

		// Check header style uses terminal colors
		headerStyle := m.headerStyle

		// We can't directly inspect the colors in lipgloss styles,
		// but we can test that the styles render without errors
		testText := "Test Header"
		rendered := headerStyle.Render(testText)

		if rendered == "" {
			t.Error("Header style should render non-empty text")
		}

		// Test footer style
		footerStyle := m.footerStyle
		footerRendered := footerStyle.Render("Test Footer")

		if footerRendered == "" {
			t.Error("Footer style should render non-empty text")
		}

		// Test that title style uses colors
		titleStyle := m.titleStyle
		titleRendered := titleStyle.Render("Test Title")

		if titleRendered == "" {
			t.Error("Title style should render non-empty text")
		}
	})

	t.Run("color_adaptation_environments", func(t *testing.T) {
		// Test color behavior in different terminal environments
		testCases := []struct {
			name     string
			termVar  string
			expected string // What we expect the terminal to support
		}{
			{"xterm_256color", "xterm-256color", "should support 256 colors"},
			{"xterm_truecolor", "xterm-256color", "should support colors"},
			{"screen", "screen", "should support basic colors"},
			{"dumb", "dumb", "should work with minimal colors"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Save original TERM
				originalTerm := os.Getenv("TERM")
				defer os.Setenv("TERM", originalTerm)

				// Set test TERM
				os.Setenv("TERM", tc.termVar)

				// Test that getTerminalThemeInfo still works
				themeInfo := getTerminalThemeInfo()

				if themeInfo == "" {
					t.Errorf("Theme info should not be empty for TERM=%s", tc.termVar)
				}

				// Should still contain "Terminal Adaptive"
				if !strings.Contains(themeInfo, "Terminal Adaptive") {
					t.Errorf("Theme info should contain 'Terminal Adaptive' for TERM=%s, got: %s", tc.termVar, themeInfo)
				}
			})
		}
	})

	t.Run("color_consistency", func(t *testing.T) {
		// Test that colors are consistent across all rendering functions
		config := &AppConfig{
			IsTestMode: true,
		}
		services := NewAppServices(config)
		m := NewModelWithServices(services)
		m.width = 100
		m.height = 30
		m.ready = true

		// Test all rendering methods to ensure they use consistent colors
		renderMethods := []struct {
			name   string
			render func() string
		}{
			{"header", func() string { return m.renderHeader() }},
			{"footer", func() string { return m.renderFooter() }},
			{"dashboard", func() string { return m.renderDashboard(20) }},
			{"settings", func() string { return m.renderSettings(20) }},
			{"help", func() string { return m.renderHelp(20) }},
		}

		for _, method := range renderMethods {
			t.Run(method.name, func(t *testing.T) {
				rendered := method.render()

				if rendered == "" {
					t.Errorf("Render method %s should return non-empty string", method.name)
				}

				// Test that rendered output doesn't contain hardcoded color codes
				// (This would indicate we're not using terminal-adaptive colors)
				hardcodedColors := []string{
					"#1a1b26", // Tokyo Night background
					"#7aa2f7", // Tokyo Night blue
					"#9ece6a", // Tokyo Night green
					"rgb(",    // RGB color functions
				}

				for _, hardcoded := range hardcodedColors {
					if strings.Contains(rendered, hardcoded) {
						t.Errorf("Render method %s contains hardcoded color %s, should use terminal-adaptive colors", method.name, hardcoded)
					}
				}
			})
		}
	})

	t.Run("termenv_integration", func(t *testing.T) {
		// Test that termenv integration works correctly

		// Test that we can create a termenv output
		termOutput := termenv.NewOutput(os.Stdout)
		if termOutput == nil {
			t.Error("Should be able to create termenv output")
			return
		}

		// Test color profile detection
		profile := termOutput.Profile
		validProfiles := []termenv.Profile{
			termenv.Ascii,
			termenv.ANSI,
			termenv.ANSI256,
			termenv.TrueColor,
		}

		profileValid := false
		for _, validProfile := range validProfiles {
			if profile == validProfile {
				profileValid = true
				break
			}
		}

		if !profileValid {
			t.Errorf("termenv should return a valid color profile, got: %v", profile)
		}

		// Test dark background detection (should not panic)
		isDark := termOutput.HasDarkBackground()
		_ = isDark // We don't assert the value since it depends on the actual terminal

		// Test color conversion
		color := termOutput.Color("4") // ANSI blue
		if color == nil {
			t.Error("termenv should be able to convert ANSI color codes")
		}
	})
}

// TestColorThemeAdaptation tests color adaptation across different themes
func TestColorThemeAdaptation(t *testing.T) {
	t.Run("color_theme_scenarios", func(t *testing.T) {
		// Simulate different terminal theme scenarios
		scenarios := []struct {
			name        string
			colorterm   string
			term        string
			description string
		}{
			{"catppuccin_alacritty", "truecolor", "xterm-256color", "Catppuccin theme in Alacritty"},
			{"tokyo_night_terminal", "truecolor", "xterm-256color", "Tokyo Night theme in terminal"},
			{"github_light_theme", "truecolor", "xterm-256color", "GitHub Light theme"},
			{"basic_terminal", "", "xterm", "Basic terminal with limited colors"},
			{"minimal_environment", "", "dumb", "Minimal terminal environment"},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				// Save original environment
				originalColorterm := os.Getenv("COLORTERM")
				originalTerm := os.Getenv("TERM")
				defer func() {
					os.Setenv("COLORTERM", originalColorterm)
					os.Setenv("TERM", originalTerm)
				}()

				// Set test environment
				os.Setenv("COLORTERM", scenario.colorterm)
				os.Setenv("TERM", scenario.term)

				// Test that our application adapts correctly
				config := &AppConfig{
					IsTestMode: true,
				}
				services := NewAppServices(config)
				m := NewModelWithServices(services)
				m.width = 100
				m.height = 30
				m.ready = true

				// Test that all views render correctly in this environment
				views := []string{
					m.renderDashboard(20),
					m.renderSettings(20),
					m.renderHelp(20),
				}

				for i, view := range views {
					if view == "" {
						t.Errorf("View %d should render in %s environment", i, scenario.name)
					}
				}

				// Test theme info reflects the environment
				themeInfo := getTerminalThemeInfo()
				if !strings.Contains(themeInfo, "Terminal Adaptive") {
					t.Errorf("Theme info should indicate terminal adaptation in %s", scenario.name)
				}
			})
		}
	})

	t.Run("color_fallback_behavior", func(t *testing.T) {
		// Test that colors fall back gracefully in limited environments

		// Save original environment
		originalColorterm := os.Getenv("COLORTERM")
		originalTerm := os.Getenv("TERM")
		defer func() {
			os.Setenv("COLORTERM", originalColorterm)
			os.Setenv("TERM", originalTerm)
		}()

		// Test with very limited terminal
		os.Setenv("COLORTERM", "")
		os.Setenv("TERM", "dumb")

		// Colors should still be defined and usable
		colors := []lipgloss.Color{
			colorBlue, colorGreen, colorRed, colorYellow,
			colorCyan, colorMagenta, colorComment, colorSubtle,
		}

		for i, color := range colors {
			if color == "" && i > 1 { // Background colors can be empty (terminal default)
				// Most colors should have some value
				colorStr := string(color)
				if len(colorStr) == 0 && i > 2 { // Allow some colors to be empty for fallback
					t.Errorf("Color %d should have fallback value in limited terminal", i)
				}
			}
		}

		// Application should still render
		config := &AppConfig{
			IsTestMode: true,
		}
		services := NewAppServices(config)
		m := NewModelWithServices(services)
		m.width = 100
		m.height = 30
		m.ready = true

		dashboard := m.renderDashboard(20)
		if dashboard == "" {
			t.Error("Dashboard should render even in limited terminal environment")
		}
	})
}

// TestTerminalThemeInfo tests the terminal theme detection function
func TestTerminalThemeInfo(t *testing.T) {
	t.Run("theme_info_format", func(t *testing.T) {
		info := getTerminalThemeInfo()

		// Should have expected format: "Terminal Adaptive (Profile, Theme)"
		if !strings.Contains(info, "Terminal Adaptive") {
			t.Errorf("Theme info should contain 'Terminal Adaptive', got: %s", info)
		}

		if !strings.Contains(info, "(") || !strings.Contains(info, ")") {
			t.Errorf("Theme info should contain parentheses with details, got: %s", info)
		}

		if !strings.Contains(info, ",") {
			t.Errorf("Theme info should contain comma separating profile and theme, got: %s", info)
		}
	})

	t.Run("theme_info_consistency", func(t *testing.T) {
		// Call multiple times to ensure consistency
		info1 := getTerminalThemeInfo()
		info2 := getTerminalThemeInfo()

		if info1 != info2 {
			t.Errorf("Theme info should be consistent across calls, got: %s vs %s", info1, info2)
		}
	})

	t.Run("theme_info_components", func(t *testing.T) {
		info := getTerminalThemeInfo()

		// Should contain color profile information
		profiles := []string{"TrueColor", "256 Color", "16 Color", "Basic"}
		hasProfile := false
		for _, profile := range profiles {
			if strings.Contains(info, profile) {
				hasProfile = true
				break
			}
		}
		if !hasProfile {
			t.Errorf("Theme info should contain color profile information, got: %s", info)
		}

		// Should contain theme type
		themes := []string{"Dark", "Light"}
		hasTheme := false
		for _, theme := range themes {
			if strings.Contains(info, theme) {
				hasTheme = true
				break
			}
		}
		if !hasTheme {
			t.Errorf("Theme info should contain theme type, got: %s", info)
		}
	})
}
