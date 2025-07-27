package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/ryanyogan/omarchy-monitor-settings/pkg/utils"
)

type MonitorDetector struct{}

func NewMonitorDetector() *MonitorDetector {
	return &MonitorDetector{}
}

func (md *MonitorDetector) DetectMonitors() ([]Monitor, error) {
	detectionMethods := []func() ([]Monitor, error){
		md.detectWithHyprctl,
		md.detectWithWlrRandr,
		md.getFallbackMonitors,
	}

	if debugMode {
		fmt.Printf("DEBUG: Starting monitor detection...\n")
	}

	for i, method := range detectionMethods {
		if debugMode {
			methodNames := []string{"hyprctl", "wlr-randr", "fallback"}
			fmt.Printf("DEBUG: Trying method %d: %s\n", i+1, methodNames[i])
		}

		monitors, err := method()
		if debugMode {
			fmt.Printf("DEBUG: Method returned %d monitors, error: %v\n", len(monitors), err)
		}

		if err == nil && len(monitors) > 0 {
			if debugMode {
				fmt.Printf("DEBUG: Successfully detected %d monitors using %s\n", len(monitors), []string{"hyprctl", "wlr-randr", "fallback"}[i])
			}
			return monitors, nil
		}
	}

	return md.getFallbackMonitors()
}

func (md *MonitorDetector) detectWithHyprctl() ([]Monitor, error) {
	if !md.commandExists("hyprctl") {
		if debugMode {
			fmt.Printf("DEBUG: hyprctl command not found\n")
		}
		return nil, fmt.Errorf("hyprctl not found")
	}

	if debugMode {
		fmt.Printf("DEBUG: Found hyprctl, running 'hyprctl monitors'\n")
	}

	cmd := exec.Command("hyprctl", "monitors")
	output, err := cmd.Output()
	if err != nil {
		if debugMode {
			fmt.Printf("DEBUG: hyprctl command failed: %v\n", err)
		}
		return nil, fmt.Errorf("failed to run hyprctl: %w", err)
	}

	if debugMode {
		fmt.Printf("DEBUG: hyprctl output (%d bytes):\n%s\n", len(output), string(output))
	}

	monitors, parseErr := md.parseHyprctlOutput(string(output))
	if debugMode {
		fmt.Printf("DEBUG: Parsed %d monitors from hyprctl output, parse error: %v\n", len(monitors), parseErr)
		for i, monitor := range monitors {
			fmt.Printf("DEBUG: Monitor %d: %s (%dx%d@%.1fHz, scale %.1f)\n",
				i, monitor.Name, monitor.Width, monitor.Height, monitor.RefreshRate, monitor.Scale)
		}
	}

	return monitors, parseErr
}

func (md *MonitorDetector) detectWithWlrRandr() ([]Monitor, error) {
	if !md.commandExists("wlr-randr") {
		return nil, fmt.Errorf("wlr-randr not found")
	}

	cmd := exec.Command("wlr-randr")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run wlr-randr: %w", err)
	}

	return md.parseWlrRandrOutput(string(output))
}

func (md *MonitorDetector) getFallbackMonitors() ([]Monitor, error) {
	var monitors []Monitor

	if runtime.GOOS == "darwin" {
		monitors = []Monitor{
			{
				Name:        "Built-in Retina Display",
				Width:       2880,
				Height:      1800,
				RefreshRate: 60.0,
				Scale:       2.0,
				Position:    Position{0, 0},
				Make:        "Apple",
				Model:       "MacBook Pro 14\"",
				IsActive:    true,
				IsPrimary:   true,
			},
		}
	} else {
		monitors = []Monitor{
			{
				Name:        "eDP-1",
				Width:       2880,
				Height:      1920,
				RefreshRate: 120.0,
				Scale:       2.0,
				Position:    Position{0, 0},
				Make:        "Framework",
				Model:       "13 Inch Laptop",
				IsActive:    true,
				IsPrimary:   true,
			},
			{
				Name:        "DP-1",
				Width:       3840,
				Height:      2160,
				RefreshRate: 60.0,
				Scale:       1.5,
				Position:    Position{2880, 0},
				Make:        "LG",
				Model:       "27UP850-W",
				IsActive:    true,
				IsPrimary:   false,
			},
		}
	}

	return monitors, nil
}

func (md *MonitorDetector) parseHyprctlOutput(output string) ([]Monitor, error) {
	var monitors []Monitor
	lines := strings.Split(output, "\n")

	var currentMonitor Monitor
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Monitor ") {
			if currentMonitor.Name != "" {
				monitors = append(monitors, currentMonitor)
			}
			currentMonitor = Monitor{}

			re := regexp.MustCompile(`Monitor ([^\s]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentMonitor.Name = matches[1]
				if debugMode {
					fmt.Printf("DEBUG: Found monitor: %s\n", currentMonitor.Name)
				}
			}
		}

		if utils.ContainsAll(line, "x", "@", " at ") {
			re := regexp.MustCompile(`(\d+)x(\d+)@([\d.]+)\s+at\s+(\d+)x(\d+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 5 {
				currentMonitor.Width = utils.ExtractInt(matches[1], 0)
				currentMonitor.Height = utils.ExtractInt(matches[2], 0)
				currentMonitor.RefreshRate = utils.ExtractFloat64(matches[3], 60.0)
				currentMonitor.Position.X = utils.ExtractInt(matches[4], 0)
				currentMonitor.Position.Y = utils.ExtractInt(matches[5], 0)
				currentMonitor.IsActive = true

				if debugMode {
					fmt.Printf("DEBUG: Parsed resolution: %dx%d@%.2fHz at %dx%d\n",
						currentMonitor.Width, currentMonitor.Height, currentMonitor.RefreshRate,
						currentMonitor.Position.X, currentMonitor.Position.Y)
				}
			} else if debugMode {
				fmt.Printf("DEBUG: Resolution regex didn't match line: '%s'\n", line)
			}
		}

		if strings.Contains(line, "scale:") {
			re := regexp.MustCompile(`scale:\s*([\d.]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentMonitor.Scale = utils.ExtractFloat64(matches[1], 1.0)
				if debugMode {
					fmt.Printf("DEBUG: Parsed scale: %.2f from line: '%s'\n", currentMonitor.Scale, line)
				}
			} else if debugMode {
				fmt.Printf("DEBUG: Scale regex didn't match line: '%s'\n", line)
			}
		}

		if strings.HasPrefix(line, "make:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				currentMonitor.Make = strings.TrimSpace(parts[1])
				if debugMode {
					fmt.Printf("DEBUG: Parsed make: '%s'\n", currentMonitor.Make)
				}
			}
		}

		if strings.HasPrefix(line, "model:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				currentMonitor.Model = strings.TrimSpace(parts[1])
				if debugMode {
					fmt.Printf("DEBUG: Parsed model: '%s'\n", currentMonitor.Model)
				}
			}
		}
	}

	if currentMonitor.Name != "" {
		monitors = append(monitors, currentMonitor)
	}

	if len(monitors) > 0 {
		monitors[0].IsPrimary = true
	}

	return monitors, nil
}

func (md *MonitorDetector) parseWlrRandrOutput(output string) ([]Monitor, error) {
	var monitors []Monitor
	lines := strings.Split(output, "\n")

	var currentMonitor Monitor
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, " ") && strings.Contains(line, " ") {
			if currentMonitor.Name != "" {
				monitors = append(monitors, currentMonitor)
			}
			currentMonitor = Monitor{}

			fields := strings.Fields(line)
			if len(fields) > 0 {
				currentMonitor.Name = fields[0]
			}
		}

		if utils.ContainsAny(line, "Hz", "*") && strings.Contains(line, "x") {
			re := regexp.MustCompile(`(\d+)x(\d+).*?([\d.]+)\s*Hz`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 3 {
				currentMonitor.Width = utils.ExtractInt(matches[1], 0)
				currentMonitor.Height = utils.ExtractInt(matches[2], 0)
				currentMonitor.RefreshRate = utils.ExtractFloat64(matches[3], 60.0)
				currentMonitor.IsActive = strings.Contains(line, "*")
				currentMonitor.Scale = 1.0
			}
		}
	}

	if currentMonitor.Name != "" {
		monitors = append(monitors, currentMonitor)
	}

	for i := range monitors {
		if monitors[i].IsActive {
			monitors[i].IsPrimary = true
			break
		}
	}

	return monitors, nil
}

func (md *MonitorDetector) commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

type ScalingManager struct{}

func NewScalingManager() *ScalingManager {
	return &ScalingManager{}
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
			{
				MonitorScale:    1.25,
				GTKScale:        1,
				FontDPI:         120,
				FontScale:       1.25,
				DisplayName:     "1.25x Maximum Space",
				Description:     "Maximum usable space with slight text scaling",
				Reasoning:       "For users who want maximum screen real estate with minimal scaling.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.25),
				EffectiveHeight: int(float64(baseHeight) / 1.25),
			},
			{
				MonitorScale:    1.0,
				GTKScale:        1,
				FontDPI:         96,
				FontScale:       1.0,
				DisplayName:     "1x Native",
				Description:     "No scaling - tiny but sharp",
				Reasoning:       "Only for users with excellent eyesight or very large monitors.",
				IsRecommended:   false,
				EffectiveWidth:  baseWidth,
				EffectiveHeight: baseHeight,
			},
		}

	case pixelCount >= 5000000:
		options = []ScalingOption{
			{
				MonitorScale:    2.0,
				GTKScale:        2,
				FontDPI:         192,
				FontScale:       1.0,
				DisplayName:     "2x Recommended",
				Description:     "Perfect for high-DPI laptop displays",
				Reasoning:       "Ideal for 13-14\" laptops with 3K displays like Framework 13. Crisp and comfortable.",
				IsRecommended:   true,
				EffectiveWidth:  baseWidth / 2,
				EffectiveHeight: baseHeight / 2,
			},
			{
				MonitorScale:    1.66667,
				GTKScale:        2,
				FontDPI:         160,
				FontScale:       1.0,
				DisplayName:     "1.67x More Space",
				Description:     "More screen space while staying readable",
				Reasoning:       "Clean 5:3 ratio scaling for productivity work on high-DPI displays.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.66667),
				EffectiveHeight: int(float64(baseHeight) / 1.66667),
			},
			{
				MonitorScale:    1.33333,
				GTKScale:        1,
				FontDPI:         128,
				FontScale:       1.33,
				DisplayName:     "1.33x Maximum Space",
				Description:     "Maximum usable space with clean scaling",
				Reasoning:       "Clean 4:3 ratio scaling for maximum productivity on high-DPI displays.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.33333),
				EffectiveHeight: int(float64(baseHeight) / 1.33333),
			},
			{
				MonitorScale:    1.0,
				GTKScale:        1,
				FontDPI:         96,
				FontScale:       1.0,
				DisplayName:     "1x Native (Tiny)",
				Description:     "No scaling - everything will be extremely small",
				Reasoning:       "Not recommended for laptop displays. Text will be barely readable.",
				IsRecommended:   false,
				EffectiveWidth:  baseWidth,
				EffectiveHeight: baseHeight,
			},
		}

	case pixelCount >= 3686400:
		options = []ScalingOption{
			{
				MonitorScale:    1.0,
				GTKScale:        1,
				FontDPI:         96,
				FontScale:       1.0,
				DisplayName:     "1x Perfect",
				Description:     "Ideal 1440p experience",
				Reasoning:       "Sweet spot for 1440p - good balance of sharpness and usability.",
				IsRecommended:   true,
				EffectiveWidth:  baseWidth,
				EffectiveHeight: baseHeight,
			},
			{
				MonitorScale:    1.25,
				GTKScale:        1,
				FontDPI:         120,
				FontScale:       1.25,
				DisplayName:     "1.25x Comfortable",
				Description:     "Slightly larger for comfort",
				Reasoning:       "Good for users who find 1x scaling slightly too small.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.25),
				EffectiveHeight: int(float64(baseHeight) / 1.25),
			},
			{
				MonitorScale:    1.5,
				GTKScale:        1,
				FontDPI:         144,
				FontScale:       1.5,
				DisplayName:     "1.5x Large",
				Description:     "Larger UI for accessibility",
				Reasoning:       "For users who prefer larger interface elements.",
				IsRecommended:   false,
				EffectiveWidth:  int(float64(baseWidth) / 1.5),
				EffectiveHeight: int(float64(baseHeight) / 1.5),
			},
		}

	case pixelCount >= 2073600:
		options = []ScalingOption{
			{
				MonitorScale:    1.0,
				GTKScale:        1,
				FontDPI:         96,
				FontScale:       1.0,
				DisplayName:     "1x Perfect",
				Description:     "Standard 1080p experience",
				Reasoning:       "Classic 1080p setup - no scaling needed for most users.",
				IsRecommended:   true,
				EffectiveWidth:  baseWidth,
				EffectiveHeight: baseHeight,
			},
			{
				MonitorScale:    1.25,
				GTKScale:        1,
				FontDPI:         120,
				FontScale:       1.25,
				DisplayName:     "1.25x Comfortable",
				Description:     "Slightly larger for small screens",
				Reasoning:       "Useful for smaller 1080p displays (13-15 inch laptops).",
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
				Description:     "No scaling applied",
				Reasoning:       "Standard scaling for this resolution.",
				IsRecommended:   true,
				EffectiveWidth:  baseWidth,
				EffectiveHeight: baseHeight,
			},
		}
	}

	return options
}

func (sm *ScalingManager) GetRecommendedScale(monitor Monitor) ScalingRecommendation {
	options := sm.GetIntelligentScalingOptions(monitor)

	for _, option := range options {
		if option.IsRecommended {
			return ScalingRecommendation{
				MonitorScale:    option.MonitorScale,
				FontScale:       option.FontScale,
				EffectiveWidth:  option.EffectiveWidth,
				EffectiveHeight: option.EffectiveHeight,
				Reasoning:       option.Reasoning,
			}
		}
	}

	if len(options) > 0 {
		first := options[0]
		return ScalingRecommendation{
			MonitorScale:    first.MonitorScale,
			FontScale:       first.FontScale,
			EffectiveWidth:  first.EffectiveWidth,
			EffectiveHeight: first.EffectiveHeight,
			Reasoning:       first.Reasoning,
		}
	}

	return ScalingRecommendation{
		MonitorScale:    1.0,
		FontScale:       1.0,
		EffectiveWidth:  monitor.Width,
		EffectiveHeight: monitor.Height,
		Reasoning:       "Default scaling",
	}
}

type ScalingRecommendation struct {
	MonitorScale    float64
	FontScale       float64
	EffectiveWidth  int
	EffectiveHeight int
	Reasoning       string
}

type ConfigManager struct {
	isDemoMode bool
}

func NewConfigManager(demoMode bool) *ConfigManager {
	return &ConfigManager{
		isDemoMode: demoMode,
	}
}

func (cm *ConfigManager) ApplyMonitorScale(monitor Monitor, scale float64) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply monitor scale %.2fx to %s\n", scale, monitor.Name)
		return nil
	}

	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		return fmt.Errorf("not running in Hyprland environment")
	}

	validatedScale := cm.validateHyprlandScale(scale)
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
		return fmt.Errorf("failed to apply scaling: Hyprland rejected scale %.3f", validatedScale)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "invalid scale") || strings.Contains(outputStr, "failed to find clean divisor") {
		return fmt.Errorf("hyprland rejected scale %.3f - try a different scaling value", validatedScale)
	}

	return nil
}

func (cm *ConfigManager) validateHyprlandScale(scale float64) float64 {
	hyprlandScales := []float64{
		1.0,
		1.25,
		1.33333,
		1.5,
		1.66667,
		1.75,
		2.0,
		2.25,
		2.5,
		3.0,
	}

	closest := hyprlandScales[0]
	minDiff := math.Abs(scale - closest)

	for _, validScale := range hyprlandScales {
		diff := math.Abs(scale - validScale)
		if diff < minDiff {
			minDiff = diff
			closest = validScale
		}
	}

	return closest
}

func (cm *ConfigManager) ApplyGTKScale(scale int) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply GTK scale %dx system-wide\n", scale)
		return nil
	}

	if err := os.Setenv("GDK_SCALE", fmt.Sprintf("%d", scale)); err != nil {
		return fmt.Errorf("failed to set GDK_SCALE: %w", err)
	}

	return nil
}

func (cm *ConfigManager) ApplyFontDPI(dpi int) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would set Xft.dpi to %d in ~/.Xresources\n", dpi)
		return nil
	}

	if err := os.Setenv("XFT_DPI", fmt.Sprintf("%d", dpi)); err != nil {
		return fmt.Errorf("failed to set XFT_DPI: %w", err)
	}

	return nil
}

func (cm *ConfigManager) ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply complete scaling option '%s':\n", option.DisplayName)
		fmt.Printf("  - Monitor scale: %.2fx\n", option.MonitorScale)
		fmt.Printf("  - GTK scale: %dx\n", option.GTKScale)
		fmt.Printf("  - Font DPI: %d\n", option.FontDPI)
		fmt.Printf("  - Additional font scale: %.2fx\n", option.FontScale)
		return nil
	}

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
		"monitor": "Changes the size of all UI elements rendered by the compositor. Applied immediately, affects everything.",
		"gtk":     "Scales GTK applications (most Linux apps). Requires logout/login. Only supports integer values (1x, 2x, 3x).",
		"font":    "Changes text size system-wide via DPI. Requires reload of applications. Fine-grained control over text rendering.",
	}
}
