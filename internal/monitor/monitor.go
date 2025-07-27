package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ryanyogan/omarchy-monitor-settings/pkg/types"
	"github.com/ryanyogan/omarchy-monitor-settings/pkg/utils"
)

type Monitor struct {
	Name        string
	Width       int
	Height      int
	RefreshRate float64
	Scale       float64
	Position    Position
	Make        string
	Model       string
	IsActive    bool
	IsPrimary   bool
}

type Position struct {
	X, Y int
}

type ScalingOption struct {
	MonitorScale    float64
	GTKScale        int
	FontDPI         int
	FontScale       float64
	DisplayName     string
	Description     string
	Reasoning       string
	IsRecommended   bool
	EffectiveWidth  int
	EffectiveHeight int
}

type DetectorInterface interface {
	DetectMonitors() ([]Monitor, error)
}

type ScalingManagerInterface interface {
	GetIntelligentScalingOptions(monitor Monitor) []ScalingOption
}

type ConfigManagerInterface interface {
	ApplyMonitorScale(monitor Monitor, scale float64) error
	ApplyGTKScale(scale int) error
	ApplyFontDPI(dpi int) error
	ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error
}

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (md *Detector) DetectMonitors() ([]Monitor, error) {
	if md.commandExists("hyprctl") {
		return md.parseHyprctlOutput()
	}

	if md.commandExists("wlr-randr") {
		return md.parseWlrRandrOutput()
	}

	return md.GetFallbackMonitors(), nil
}

func (md *Detector) parseHyprctlOutput() ([]Monitor, error) {
	cmd := exec.Command("hyprctl", "monitors")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var monitors []Monitor
	lines := strings.Split(string(output), "\n")
	var currentMonitor *Monitor

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Monitor ") {
			if currentMonitor != nil {
				monitors = append(monitors, *currentMonitor)
			}

			nameMatch := regexp.MustCompile(`Monitor (\S+)`).FindStringSubmatch(line)
			if len(nameMatch) > 1 {
				currentMonitor = &Monitor{
					Name: nameMatch[1],
				}
			}
			continue
		}

		if currentMonitor == nil {
			continue
		}

		if strings.Contains(line, "x") && strings.Contains(line, "@") {
			resolutionMatch := regexp.MustCompile(`(\d+)x(\d+)@([\d.]+)Hz`).FindStringSubmatch(line)
			if len(resolutionMatch) > 3 {
				if width, err := strconv.Atoi(resolutionMatch[1]); err == nil {
					currentMonitor.Width = width
				}
				if height, err := strconv.Atoi(resolutionMatch[2]); err == nil {
					currentMonitor.Height = height
				}
				if refreshRate, err := strconv.ParseFloat(resolutionMatch[3], 64); err == nil {
					currentMonitor.RefreshRate = refreshRate
				}
			}
		}

		if strings.HasPrefix(line, "scale:") {
			scaleMatch := regexp.MustCompile(`scale:\s*([\d.]+)`).FindStringSubmatch(line)
			if len(scaleMatch) > 1 {
				if scale, err := strconv.ParseFloat(scaleMatch[1], 64); err == nil {
					currentMonitor.Scale = scale
				}
			}
		}

		if strings.HasPrefix(line, "description:") {
			descMatch := regexp.MustCompile(`description:\s*(.+)`).FindStringSubmatch(line)
			if len(descMatch) > 1 {
				currentMonitor.Model = strings.TrimSpace(descMatch[1])
			}
		}

		if strings.HasPrefix(line, "make:") {
			makeMatch := regexp.MustCompile(`make:\s*(.+)`).FindStringSubmatch(line)
			if len(makeMatch) > 1 {
				currentMonitor.Make = strings.TrimSpace(makeMatch[1])
			}
		}

		if strings.Contains(line, "focused: yes") {
			currentMonitor.IsActive = true
			currentMonitor.IsPrimary = true
		}
	}

	if currentMonitor != nil {
		monitors = append(monitors, *currentMonitor)
	}

	return monitors, nil
}

func (md *Detector) parseWlrRandrOutput() ([]Monitor, error) {
	cmd := exec.Command("wlr-randr")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute wlr-randr: %w", err)
	}

	var monitors []Monitor
	lines := strings.Split(string(output), "\n")
	var currentMonitor *Monitor

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, " ") && strings.Contains(line, " ") {
			if currentMonitor != nil {
				monitors = append(monitors, *currentMonitor)
			}

			parts := strings.Fields(line)
			if len(parts) > 0 {
				currentMonitor = &Monitor{
					Name: parts[0],
				}
			}
			continue
		}

		if currentMonitor == nil {
			continue
		}

		if strings.Contains(line, "x") && strings.Contains(line, "@") {
			resolutionMatch := regexp.MustCompile(`(\d+)x(\d+)@([\d.]+)Hz`).FindStringSubmatch(line)
			if len(resolutionMatch) > 3 {
				if width, err := strconv.Atoi(resolutionMatch[1]); err == nil {
					currentMonitor.Width = width
				}
				if height, err := strconv.Atoi(resolutionMatch[2]); err == nil {
					currentMonitor.Height = height
				}
				if refreshRate, err := strconv.ParseFloat(resolutionMatch[3], 64); err == nil {
					currentMonitor.RefreshRate = refreshRate
				}
			}
		}

		if strings.Contains(line, "scale:") {
			scaleMatch := regexp.MustCompile(`scale:\s*([\d.]+)`).FindStringSubmatch(line)
			if len(scaleMatch) > 1 {
				if scale, err := strconv.ParseFloat(scaleMatch[1], 64); err == nil {
					currentMonitor.Scale = scale
				}
			}
		}
	}

	if currentMonitor != nil {
		monitors = append(monitors, *currentMonitor)
	}

	return monitors, nil
}

func (md *Detector) GetFallbackMonitors() []Monitor {
	return []Monitor{
		{
			Name:        "eDP-1",
			Width:       1920,
			Height:      1080,
			RefreshRate: 60.0,
			Scale:       1.0,
			Position:    Position{X: 0, Y: 0},
			Make:        "Demo",
			Model:       "Display",
			IsActive:    true,
			IsPrimary:   true,
		},
		{
			Name:        "HDMI-A-1",
			Width:       2560,
			Height:      1440,
			RefreshRate: 60.0,
			Scale:       1.0,
			Position:    Position{X: 1920, Y: 0},
			Make:        "Demo",
			Model:       "External",
			IsActive:    true,
			IsPrimary:   false,
		},
	}
}

func (md *Detector) commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

type ScalingManager struct{}

func NewScalingManager() *ScalingManager {
	return &ScalingManager{}
}

func (sm *ScalingManager) GetIntelligentScalingOptions(monitor Monitor) []ScalingOption {
	pixelCount := monitor.Width * monitor.Height

	var options []ScalingOption

	baseWidth := monitor.Width
	baseHeight := monitor.Height

	switch {
	case pixelCount >= 8294400:
		options = []ScalingOption{
			{
				MonitorScale:    2.0,
				GTKScale:        2,
				FontDPI:         192,
				FontScale:       1.0,
				DisplayName:     "2x Perfect",
				Description:     "Sharp 4K experience with crisp text",
				Reasoning:       "Industry standard for 4K displays. Perfect integer scaling with no blur.",
				IsRecommended:   true,
				EffectiveWidth:  baseWidth / 2,
				EffectiveHeight: baseHeight / 2,
			},
			{
				MonitorScale:    1.5,
				GTKScale:        1,
				FontDPI:         144,
				FontScale:       1.5,
				DisplayName:     "1.5x Balanced",
				Description:     "More screen space with readable text",
				Reasoning:       "Good compromise between space and readability for productivity.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.5),
				EffectiveHeight: int(float64(baseHeight) / 1.5),
			},
		}
	case pixelCount >= 3686400:
		options = []ScalingOption{
			{
				MonitorScale:    1.5,
				GTKScale:        1,
				FontDPI:         144,
				FontScale:       1.5,
				DisplayName:     "1.5x Sharp",
				Description:     "Perfect scaling for 1440p displays",
				Reasoning:       "Ideal for 1440p displays. Provides crisp text and good screen real estate.",
				IsRecommended:   true,
				EffectiveWidth:  int(float64(baseWidth) / 1.5),
				EffectiveHeight: int(float64(baseHeight) / 1.5),
			},
			{
				MonitorScale:    1.25,
				GTKScale:        1,
				FontDPI:         120,
				FontScale:       1.25,
				DisplayName:     "1.25x Balanced",
				Description:     "More space with readable text",
				Reasoning:       "Good balance between space and readability for productivity work.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.25),
				EffectiveHeight: int(float64(baseHeight) / 1.25),
			},
		}
	default:
		options = []ScalingOption{
			{
				MonitorScale:    1.0,
				GTKScale:        1,
				FontDPI:         96,
				FontScale:       1.0,
				DisplayName:     "1x Native",
				Description:     "Native resolution with standard scaling",
				Reasoning:       "Standard scaling for 1080p and lower displays. Good for most use cases.",
				IsRecommended:   true,
				EffectiveWidth:  baseWidth,
				EffectiveHeight: baseHeight,
			},
			{
				MonitorScale:    1.25,
				GTKScale:        1,
				FontDPI:         120,
				FontScale:       1.25,
				DisplayName:     "1.25x Enhanced",
				Description:     "Slightly larger text for better readability",
				Reasoning:       "Good for users who prefer larger text without losing too much screen space.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.25),
				EffectiveHeight: int(float64(baseHeight) / 1.25),
			},
		}
	}

	return options
}

func (sm *ScalingManager) GetRecommendedScale(monitor Monitor) float64 {
	options := sm.GetIntelligentScalingOptions(monitor)
	for _, option := range options {
		if option.IsRecommended {
			return option.MonitorScale
		}
	}
	return 1.0
}

type ScalingRecommendation struct {
	MonitorScale float64
	GTKScale     int
	FontDPI      int
	Reasoning    string
}

type ConfigManager struct {
	isDemoMode bool
}

func NewConfigManager(isDemoMode bool) *ConfigManager {
	return &ConfigManager{
		isDemoMode: isDemoMode,
	}
}

func (cm *ConfigManager) ApplyMonitorScale(monitor Monitor, scale float64) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply monitor scale %.2fx to %s\n", scale, monitor.Name)
		return nil
	}

	validatedScale := utils.ValidateMonitorScale(scale, types.MinMonitorScale, types.MaxMonitorScale)
	if validatedScale != scale {
		fmt.Printf("Adjusted scale from %.3f to %.3f for Hyprland compatibility\n", scale, validatedScale)
	}

	monitorName := strings.ReplaceAll(monitor.Name, " ", "_")
	monitorName = strings.ReplaceAll(monitorName, ";", "")
	monitorName = strings.ReplaceAll(monitorName, "&", "")
	monitorName = strings.ReplaceAll(monitorName, "|", "")
	monitorName = strings.ReplaceAll(monitorName, "`", "")
	monitorName = strings.ReplaceAll(monitorName, "$", "")
	monitorName = strings.ReplaceAll(monitorName, "(", "")
	monitorName = strings.ReplaceAll(monitorName, ")", "")
	monitorName = strings.ReplaceAll(monitorName, "'", "")
	monitorName = strings.ReplaceAll(monitorName, "\"", "")

	cmd := exec.Command("hyprctl", "keyword", "monitor", fmt.Sprintf("%s,preferred,auto,%.5f", monitorName, validatedScale)) // nosec G204

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply monitor scale: %w, output: %s", err, string(output))
	}

	return nil
}

func (cm *ConfigManager) ApplyGTKScale(scale int) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply GTK scale %dx system-wide\n", scale)
		return nil
	}

	validatedScale := utils.ValidateGTKScale(scale, types.MinGTKScale, types.MaxGTKScale)

	if err := os.Setenv("GDK_SCALE", fmt.Sprintf("%d", validatedScale)); err != nil {
		return fmt.Errorf("failed to set GDK_SCALE: %w", err)
	}

	return nil
}

func (cm *ConfigManager) ApplyFontDPI(dpi int) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would set Xft.dpi to %d in ~/.Xresources\n", dpi)
		return nil
	}

	validatedDPI := utils.ValidateFontDPI(dpi, types.MinFontDPI, types.MaxFontDPI)

	if err := os.Setenv("XFT_DPI", fmt.Sprintf("%d", validatedDPI)); err != nil {
		return fmt.Errorf("failed to set XFT_DPI: %w", err)
	}

	return nil
}

func (cm *ConfigManager) ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error {
	if err := cm.ApplyMonitorScale(monitor, option.MonitorScale); err != nil {
		return fmt.Errorf("failed to apply monitor scale: %w", err)
	}

	if err := cm.ApplyGTKScale(option.GTKScale); err != nil {
		return fmt.Errorf("failed to apply GTK scale: %w", err)
	}

	if err := cm.ApplyFontDPI(option.FontDPI); err != nil {
		return fmt.Errorf("failed to apply font DPI: %w", err)
	}

	return nil
}

func (cm *ConfigManager) GetScalingExplanations() map[string]string {
	return map[string]string{
		"monitor_scale": "Controls the compositor-level scaling. Affects the entire display output.",
		"gtk_scale":     "Controls GTK application scaling. Affects GTK-based applications.",
		"font_dpi":      "Controls font rendering DPI. Affects text size and clarity.",
	}
}
