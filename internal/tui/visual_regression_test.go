package tui

import (
	"fmt"
	"testing"

	"github.com/ryanyogan/omarchy-monitor-settings/internal/app"
	"github.com/ryanyogan/omarchy-monitor-settings/internal/monitor"
	visualtest "github.com/ryanyogan/omarchy-monitor-settings/pkg/testing"
)

func TestVisualRegression(t *testing.T) {
	vt := visualtest.NewVisualTester(t, "testdata/golden")

	screenSizes := []struct{ Width, Height int }{
		{80, 24},
		{100, 30},
		{120, 40},
		{150, 50},
		{200, 60},
	}

	t.Run("Dashboard", func(t *testing.T) {
		model := createTestModelForVisual(ModeDashboard)
		vt.MultiSizeTest("dashboard", model, screenSizes)
	})

	t.Run("MonitorSelection", func(t *testing.T) {
		model := createTestModelForVisual(ModeMonitorSelection)
		vt.MultiSizeTest("monitor_selection", model, screenSizes)
	})

	t.Run("ScalingOptions", func(t *testing.T) {
		model := createTestModelForVisual(ModeScalingOptions)
		vt.MultiSizeTest("scaling_options", model, screenSizes)
	})

	t.Run("ManualScaling", func(t *testing.T) {
		model := createTestModelForVisual(ModeManualScaling)
		vt.MultiSizeTest("manual_scaling", model, screenSizes)
	})

	t.Run("Settings", func(t *testing.T) {
		model := createTestModelForVisual(ModeSettings)
		vt.MultiSizeTest("settings", model, screenSizes)
	})

	t.Run("Help", func(t *testing.T) {
		model := createTestModelForVisual(ModeHelp)
		vt.MultiSizeTest("help", model, screenSizes)
	})
}

func TestVisualRegressionEdgeCases(t *testing.T) {
	vt := visualtest.NewVisualTester(t, "testdata/golden")

	t.Run("TerminalTooSmall", func(t *testing.T) {
		smallSizes := []struct{ Width, Height int }{
			{70, 20},
			{80, 15},
			{60, 10},
		}

		model := createTestModelForVisual(ModeDashboard)
		vt.MultiSizeTest("terminal_too_small", model, smallSizes)
	})

	t.Run("NoMonitors", func(t *testing.T) {
		model := createTestModelWithMonitors(ModeDashboard, []monitor.Monitor{})
		vt.TestVisualRegression(visualtest.VisualTestConfig{
			Name:   "no_monitors",
			Width:  120,
			Height: 40,
			Model:  model,
		})
	})

	t.Run("SingleMonitor", func(t *testing.T) {
		testMonitor := monitor.Monitor{
			Name:        "HDMI-1",
			Make:        "Samsung",
			Model:       "U2414H",
			Width:       1920,
			Height:      1080,
			RefreshRate: 60.0,
			Scale:       1.0,
			IsActive:    true,
		}

		model := createTestModelWithMonitors(ModeDashboard, []monitor.Monitor{testMonitor})
		vt.TestVisualRegression(visualtest.VisualTestConfig{
			Name:   "single_monitor",
			Width:  120,
			Height: 40,
			Model:  model,
		})
	})

	t.Run("ManyMonitors", func(t *testing.T) {
		monitors := make([]monitor.Monitor, 5)
		for i := range monitors {
			monitors[i] = monitor.Monitor{
				Name:        fmt.Sprintf("HDMI-%d", i+1),
				Make:        "Dell",
				Model:       fmt.Sprintf("U2414H-%d", i+1),
				Width:       1920,
				Height:      1080,
				RefreshRate: 60.0,
				Scale:       1.0,
				IsActive:    i == 0,
			}
		}

		model := createTestModelWithMonitors(ModeDashboard, monitors)
		vt.TestVisualRegression(visualtest.VisualTestConfig{
			Name:   "many_monitors",
			Width:  150,
			Height: 50,
			Model:  model,
		})
	})
}

func TestVisualRegressionInteractions(t *testing.T) {
	vt := visualtest.NewVisualTester(t, "testdata/golden")

	t.Run("NavigationStates", func(t *testing.T) {
		for i := 0; i < 7; i++ {
			model := createTestModelForVisual(ModeDashboard)
			model.selectedOption = i

			vt.TestVisualRegression(visualtest.VisualTestConfig{
				Name:   fmt.Sprintf("navigation_selected_%d", i),
				Width:  120,
				Height: 40,
				Model:  model,
			})
		}
	})

	t.Run("ManualScalingControls", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			model := createTestModelForVisual(ModeManualScaling)
			model.selectedManualControl = i

			vt.TestVisualRegression(visualtest.VisualTestConfig{
				Name:   fmt.Sprintf("manual_scaling_control_%d", i),
				Width:  120,
				Height: 40,
				Model:  model,
			})
		}
	})

	t.Run("ScalingValues", func(t *testing.T) {
		testCases := []struct {
			name         string
			monitorScale float64
			gtkScale     int
			fontDPI      int
		}{
			{"min_values", 0.5, 1, 72},
			{"default_values", 1.0, 1, 96},
			{"high_values", 2.0, 3, 144},
			{"max_values", 4.0, 3, 300},
		}

		for _, tc := range testCases {
			model := createTestModelForVisual(ModeManualScaling)
			model.manualMonitorScale = tc.monitorScale
			model.manualGTKScale = tc.gtkScale
			model.manualFontDPI = tc.fontDPI

			vt.TestVisualRegression(visualtest.VisualTestConfig{
				Name:   fmt.Sprintf("scaling_values_%s", tc.name),
				Width:  120,
				Height: 40,
				Model:  model,
			})
		}
	})
}

func TestVisualRegressionThemes(t *testing.T) {
	vt := visualtest.NewVisualTester(t, "testdata/golden")

	terminalEnvs := []struct {
		name      string
		term      string
		colorterm string
	}{
		{"xterm_256color", "xterm-256color", "truecolor"},
		{"xterm_basic", "xterm", ""},
		{"screen", "screen", ""},
		{"tmux", "tmux-256color", "truecolor"},
	}

	for _, env := range terminalEnvs {
		t.Run(env.name, func(t *testing.T) {
			t.Setenv("TERM", env.term)
			if env.colorterm != "" {
				t.Setenv("COLORTERM", env.colorterm)
			}

			model := createTestModelForVisual(ModeDashboard)
			vt.TestVisualRegression(visualtest.VisualTestConfig{
				Name:   fmt.Sprintf("theme_%s", env.name),
				Width:  120,
				Height: 40,
				Model:  model,
			})
		})
	}
}

func createTestModelForVisual(mode AppMode) Model {
	config := &app.Config{
		DebugMode:       false,
		NoHyprlandCheck: true,
	}

	services := &app.Services{
		Config:          config,
		MonitorDetector: &MockMonitorDetector{},
		ScalingManager:  &MockScalingManager{},
		ConfigManager:   &MockConfigManager{},
	}

	model := NewModelWithServices(services)
	model.mode = mode
	model.ready = true

	model.monitors = []monitor.Monitor{
		{
			Name:        "HDMI-A-1",
			Make:        "Dell",
			Model:       "U2414H",
			Width:       1920,
			Height:      1080,
			RefreshRate: 60.0,
			Scale:       1.0,
			IsActive:    true,
		},
		{
			Name:        "DP-1",
			Make:        "Samsung",
			Model:       "C27F390",
			Width:       1920,
			Height:      1080,
			RefreshRate: 75.0,
			Scale:       1.25,
			IsActive:    false,
		},
	}

	// Update scaling options for the first monitor
	if len(model.monitors) > 0 {
		model.scalingOptions = services.ScalingManager.GetIntelligentScalingOptions(model.monitors[0])
	}

	return model
}

func createTestModelWithMonitors(mode AppMode, monitors []monitor.Monitor) Model {
	model := createTestModelForVisual(mode)
	model.monitors = monitors
	return model
}

type MockMonitorDetector struct{}

func (m *MockMonitorDetector) DetectMonitors() ([]monitor.Monitor, error) {
	return []monitor.Monitor{
		{
			Name:        "HDMI-A-1",
			Make:        "Dell",
			Model:       "U2414H",
			Width:       1920,
			Height:      1080,
			RefreshRate: 60.0,
			Scale:       1.0,
			IsActive:    true,
		},
	}, nil
}

type MockScalingManager struct{}

func (m *MockScalingManager) GetIntelligentScalingOptions(mon monitor.Monitor) []monitor.ScalingOption {
	return []monitor.ScalingOption{
		{
			DisplayName:     "1x Native",
			Description:     "Native resolution with standard scaling",
			MonitorScale:    1.0,
			GTKScale:        1,
			FontDPI:         96,
			IsRecommended:   true,
			EffectiveWidth:  mon.Width,
			EffectiveHeight: mon.Height,
		},
		{
			DisplayName:     "1.25x Enhanced",
			Description:     "Slightly larger text for better readability",
			MonitorScale:    1.25,
			GTKScale:        1,
			FontDPI:         120,
			IsRecommended:   false,
			EffectiveWidth:  int(float64(mon.Width) / 1.25),
			EffectiveHeight: int(float64(mon.Height) / 1.25),
		},
		{
			DisplayName:     "1.5x Large",
			Description:     "Larger text for accessibility",
			MonitorScale:    1.5,
			GTKScale:        1,
			FontDPI:         144,
			IsRecommended:   false,
			EffectiveWidth:  int(float64(mon.Width) / 1.5),
			EffectiveHeight: int(float64(mon.Height) / 1.5),
		},
	}
}

type MockConfigManager struct{}

func (m *MockConfigManager) ApplyMonitorScale(monitor monitor.Monitor, scale float64) error {
	return nil
}

func (m *MockConfigManager) ApplyGTKScale(scale int) error {
	return nil
}

func (m *MockConfigManager) ApplyFontDPI(dpi int) error {
	return nil
}

func (m *MockConfigManager) ApplyCompleteScalingOption(mon monitor.Monitor, option monitor.ScalingOption) error {
	return nil
}
