package main

import (
	"fmt"
	"testing"

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
		model := createTestModel(ModeDashboard)
		vt.MultiSizeTest("dashboard", model, screenSizes)
	})

	t.Run("MonitorSelection", func(t *testing.T) {
		model := createTestModel(ModeMonitorSelection)
		vt.MultiSizeTest("monitor_selection", model, screenSizes)
	})

	t.Run("ScalingOptions", func(t *testing.T) {
		model := createTestModel(ModeScalingOptions)
		vt.MultiSizeTest("scaling_options", model, screenSizes)
	})

	t.Run("ManualScaling", func(t *testing.T) {
		model := createTestModel(ModeManualScaling)
		vt.MultiSizeTest("manual_scaling", model, screenSizes)
	})

	t.Run("Settings", func(t *testing.T) {
		model := createTestModel(ModeSettings)
		vt.MultiSizeTest("settings", model, screenSizes)
	})

	t.Run("Help", func(t *testing.T) {
		model := createTestModel(ModeHelp)
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

		model := createTestModel(ModeDashboard)
		vt.MultiSizeTest("terminal_too_small", model, smallSizes)
	})

	t.Run("NoMonitors", func(t *testing.T) {
		model := createTestModelWithMonitors(ModeDashboard, []Monitor{})
		vt.TestVisualRegression(visualtest.VisualTestConfig{
			Name:   "no_monitors",
			Width:  120,
			Height: 40,
			Model:  model,
		})
	})

	t.Run("SingleMonitor", func(t *testing.T) {
		monitor := Monitor{
			Name:        "HDMI-1",
			Make:        "Samsung",
			Model:       "U2414H",
			Width:       1920,
			Height:      1080,
			RefreshRate: 60.0,
			Scale:       1.0,
			IsActive:    true,
		}

		model := createTestModelWithMonitors(ModeDashboard, []Monitor{monitor})
		vt.TestVisualRegression(visualtest.VisualTestConfig{
			Name:   "single_monitor",
			Width:  120,
			Height: 40,
			Model:  model,
		})
	})

	t.Run("ManyMonitors", func(t *testing.T) {
		monitors := make([]Monitor, 5)
		for i := range monitors {
			monitors[i] = Monitor{
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
			model := createTestModel(ModeDashboard)
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
			model := createTestModel(ModeManualScaling)
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
			model := createTestModel(ModeManualScaling)
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

			model := createTestModel(ModeDashboard)
			vt.TestVisualRegression(visualtest.VisualTestConfig{
				Name:   fmt.Sprintf("theme_%s", env.name),
				Width:  120,
				Height: 40,
				Model:  model,
			})
		})
	}
}

func createTestModel(mode AppMode) Model {
	config := &AppConfig{
		DebugMode:       false,
		NoHyprlandCheck: true,
	}

	services := &AppServices{
		Config:          config,
		MonitorDetector: &MockMonitorDetector{},
		ScalingManager:  &MockScalingManager{},
		ConfigManager:   &MockConfigManager{},
	}

	model := NewModelWithServices(services)
	model.mode = mode
	model.ready = true

	model.monitors = []Monitor{
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

	return model
}

func createTestModelWithMonitors(mode AppMode, monitors []Monitor) Model {
	model := createTestModel(mode)
	model.monitors = monitors
	return model
}

type MockMonitorDetector struct{}

func (m *MockMonitorDetector) DetectMonitors() ([]Monitor, error) {
	return []Monitor{
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

func (m *MockScalingManager) GetIntelligentScalingOptions(monitor Monitor) []ScalingOption {
	return []ScalingOption{
		{
			DisplayName:     "Standard (1x)",
			Description:     "Default scaling for sharp text",
			MonitorScale:    1.0,
			GTKScale:        1,
			FontDPI:         96,
			IsRecommended:   true,
			EffectiveWidth:  monitor.Width,
			EffectiveHeight: monitor.Height,
		},
		{
			DisplayName:     "Large (1.25x)",
			Description:     "Comfortable scaling for most users",
			MonitorScale:    1.25,
			GTKScale:        1,
			FontDPI:         120,
			IsRecommended:   false,
			EffectiveWidth:  int(float64(monitor.Width) / 1.25),
			EffectiveHeight: int(float64(monitor.Height) / 1.25),
		},
	}
}

type MockConfigManager struct{}

func (m *MockConfigManager) ApplyMonitorScale(monitor Monitor, scale float64) error {
	return nil
}

func (m *MockConfigManager) ApplyGTKScale(scale int) error {
	return nil
}

func (m *MockConfigManager) ApplyFontDPI(dpi int) error {
	return nil
}

func (m *MockConfigManager) ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error {
	return nil
}
