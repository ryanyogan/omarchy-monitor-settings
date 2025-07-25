package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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
			expected: "Hyprland Monitor Manager",
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
