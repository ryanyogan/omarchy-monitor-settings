package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// MonitorDetector handles monitor detection across different platforms and tools
type MonitorDetector struct{}

// NewMonitorDetector creates a new monitor detector
func NewMonitorDetector() *MonitorDetector {
	return &MonitorDetector{}
}

// DetectMonitors attempts to detect monitors using various methods
func (md *MonitorDetector) DetectMonitors() ([]Monitor, error) {
	// Try different detection methods in order of preference
	detectionMethods := []func() ([]Monitor, error){
		md.detectWithHyprctl,
		md.detectWithWlrRandr,
		md.detectWithXrandr,
		md.getFallbackMonitors, // Always succeeds with demo data
	}

	if debugMode {
		fmt.Printf("DEBUG: Starting monitor detection...\n")
	}

	for i, method := range detectionMethods {
		if debugMode {
			methodNames := []string{"hyprctl", "wlr-randr", "xrandr", "fallback"}
			fmt.Printf("DEBUG: Trying method %d: %s\n", i+1, methodNames[i])
		}

		monitors, err := method()
		if debugMode {
			fmt.Printf("DEBUG: Method returned %d monitors, error: %v\n", len(monitors), err)
		}

		if err == nil && len(monitors) > 0 {
			if debugMode {
				fmt.Printf("DEBUG: Successfully detected %d monitors using %s\n", len(monitors), []string{"hyprctl", "wlr-randr", "xrandr", "fallback"}[i])
			}
			return monitors, nil
		}
	}

	// This should never be reached due to fallback, but just in case
	return md.getFallbackMonitors()
}

// detectWithHyprctl detects monitors using hyprctl (Hyprland's control utility)
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

// detectWithWlrRandr detects monitors using wlr-randr (Wayland display manager)
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

// detectWithXrandr detects monitors using xrandr (X11 display manager - fallback)
func (md *MonitorDetector) detectWithXrandr() ([]Monitor, error) {
	if !md.commandExists("xrandr") {
		return nil, fmt.Errorf("xrandr not found")
	}

	cmd := exec.Command("xrandr")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run xrandr: %w", err)
	}

	return md.parseXrandrOutput(string(output))
}

// getFallbackMonitors provides realistic demo monitors for testing
func (md *MonitorDetector) getFallbackMonitors() ([]Monitor, error) {
	// Provide different demo data based on platform for realistic testing
	var monitors []Monitor

	if runtime.GOOS == "darwin" {
		// macOS-like monitors for testing
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
		// Linux/Framework-like monitors for testing
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

// parseHyprctlOutput parses the output from hyprctl monitors
func (md *MonitorDetector) parseHyprctlOutput(output string) ([]Monitor, error) {
	var monitors []Monitor
	lines := strings.Split(output, "\n")

	var currentMonitor Monitor
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Monitor line: "Monitor eDP-1 (ID 0):"
		if strings.HasPrefix(line, "Monitor ") {
			if currentMonitor.Name != "" {
				monitors = append(monitors, currentMonitor)
			}
			currentMonitor = Monitor{}

			// Extract monitor name
			re := regexp.MustCompile(`Monitor ([^\s]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentMonitor.Name = matches[1]
				if debugMode {
					fmt.Printf("DEBUG: Found monitor: %s\n", currentMonitor.Name)
				}
			}
		}

		// Resolution line: "2880x1920@120.00000 at 0x0" (note: no Hz suffix in newer Hyprland)
		if strings.Contains(line, "x") && strings.Contains(line, "@") && strings.Contains(line, " at ") {
			re := regexp.MustCompile(`(\d+)x(\d+)@([\d.]+)\s+at\s+(\d+)x(\d+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 5 {
				currentMonitor.Width, _ = strconv.Atoi(matches[1])
				currentMonitor.Height, _ = strconv.Atoi(matches[2])
				currentMonitor.RefreshRate, _ = strconv.ParseFloat(matches[3], 64)
				currentMonitor.Position.X, _ = strconv.Atoi(matches[4])
				currentMonitor.Position.Y, _ = strconv.Atoi(matches[5])
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

		// Scale line: "scale: 1.67" or "scale: 2.00"
		if strings.Contains(line, "scale:") {
			re := regexp.MustCompile(`scale:\s*([\d.]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentMonitor.Scale, _ = strconv.ParseFloat(matches[1], 64)
				if debugMode {
					fmt.Printf("DEBUG: Parsed scale: %.2f from line: '%s'\n", currentMonitor.Scale, line)
				}
			} else if debugMode {
				fmt.Printf("DEBUG: Scale regex didn't match line: '%s'\n", line)
			}
		}

		// Make line: "make: BOE"
		if strings.HasPrefix(line, "make:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				currentMonitor.Make = strings.TrimSpace(parts[1])
				if debugMode {
					fmt.Printf("DEBUG: Parsed make: '%s'\n", currentMonitor.Make)
				}
			}
		}

		// Model line: "model: NE135A1M-NY1"
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

	// Add the last monitor
	if currentMonitor.Name != "" {
		monitors = append(monitors, currentMonitor)
	}

	// Set the first monitor as primary if none specified
	if len(monitors) > 0 {
		monitors[0].IsPrimary = true
	}

	return monitors, nil
}

// parseWlrRandrOutput parses the output from wlr-randr
func (md *MonitorDetector) parseWlrRandrOutput(output string) ([]Monitor, error) {
	var monitors []Monitor
	lines := strings.Split(output, "\n")

	var currentMonitor Monitor
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Monitor line format varies, try to detect monitor names
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

		// Look for resolution and refresh rate
		if strings.Contains(line, "x") && (strings.Contains(line, "Hz") || strings.Contains(line, "*")) {
			re := regexp.MustCompile(`(\d+)x(\d+).*?([\d.]+)\s*Hz`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 3 {
				currentMonitor.Width, _ = strconv.Atoi(matches[1])
				currentMonitor.Height, _ = strconv.Atoi(matches[2])
				currentMonitor.RefreshRate, _ = strconv.ParseFloat(matches[3], 64)
				currentMonitor.IsActive = strings.Contains(line, "*")
				currentMonitor.Scale = 1.0 // Default scale
			}
		}
	}

	// Add the last monitor
	if currentMonitor.Name != "" {
		monitors = append(monitors, currentMonitor)
	}

	// Set the first active monitor as primary
	for i := range monitors {
		if monitors[i].IsActive {
			monitors[i].IsPrimary = true
			break
		}
	}

	return monitors, nil
}

// parseXrandrOutput parses the output from xrandr
func (md *MonitorDetector) parseXrandrOutput(output string) ([]Monitor, error) {
	var monitors []Monitor
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for connected monitors: "DP1 connected 1920x1080+1366+0"
		if strings.Contains(line, " connected") {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}

			monitor := Monitor{
				Name:     fields[0],
				Scale:    1.0,
				IsActive: true,
			}

			// Parse resolution from "1920x1080+1366+0" format
			if len(fields) > 2 {
				resolutionStr := fields[2]
				re := regexp.MustCompile(`(\d+)x(\d+)(?:\+(\d+)\+(\d+))?`)
				matches := re.FindStringSubmatch(resolutionStr)
				if len(matches) > 2 {
					monitor.Width, _ = strconv.Atoi(matches[1])
					monitor.Height, _ = strconv.Atoi(matches[2])
					if len(matches) > 4 {
						monitor.Position.X, _ = strconv.Atoi(matches[3])
						monitor.Position.Y, _ = strconv.Atoi(matches[4])
					}
				}
			}

			// Look for refresh rate in subsequent lines
			monitor.RefreshRate = 60.0 // Default

			monitors = append(monitors, monitor)
		}
	}

	// Set the first monitor as primary
	if len(monitors) > 0 {
		monitors[0].IsPrimary = true
	}

	return monitors, nil
}

// commandExists checks if a command is available in PATH
func (md *MonitorDetector) commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// ScalingManager handles scaling recommendations and configuration
type ScalingManager struct{}

// NewScalingManager creates a new scaling manager
func NewScalingManager() *ScalingManager {
	return &ScalingManager{}
}

// ScalingOption represents a scaling recommendation with detailed information
type ScalingOption struct {
	MonitorScale    float64
	GTKScale        int     // GTK only supports integer scaling
	FontDPI         int     // Xft.dpi setting (base 96)
	FontScale       float64 // Additional font scaling
	DisplayName     string
	Description     string
	Reasoning       string
	IsRecommended   bool
	EffectiveWidth  int
	EffectiveHeight int
}

// GetIntelligentScalingOptions returns research-based scaling recommendations
func (sm *ScalingManager) GetIntelligentScalingOptions(monitor Monitor) []ScalingOption {
	pixelCount := monitor.Width * monitor.Height

	var options []ScalingOption

	// Base calculations for different scenarios
	baseWidth := monitor.Width
	baseHeight := monitor.Height

	switch {
	case pixelCount >= 8294400: // 4K+ (3840x2160+)
		options = []ScalingOption{
			{
				MonitorScale:    2.0,
				GTKScale:        2,
				FontDPI:         192, // 96 * 2
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
				FontDPI:         144, // 96 * 1.5
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
				FontDPI:         120, // 96 * 1.25
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

	case pixelCount >= 5000000: // High-DPI displays like Framework 13 (2880x1920), MacBook Pro 13", etc.
		options = []ScalingOption{
			{
				MonitorScale:    2.0,
				GTKScale:        2,
				FontDPI:         192, // 96 * 2
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
				FontDPI:         160, // 96 * 1.66667
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
				FontDPI:         128, // 96 * 1.33333
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

	case pixelCount >= 3686400: // 1440p (2560x1440)
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

	case pixelCount >= 2073600: // 1080p (1920x1080)
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

	default: // Lower resolutions
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

// GetRecommendedScale returns the single best scaling recommendation (legacy function)
func (sm *ScalingManager) GetRecommendedScale(monitor Monitor) ScalingRecommendation {
	options := sm.GetIntelligentScalingOptions(monitor)

	// Find the recommended option
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

	// Fallback to first option
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

	// Ultimate fallback
	return ScalingRecommendation{
		MonitorScale:    1.0,
		FontScale:       1.0,
		EffectiveWidth:  monitor.Width,
		EffectiveHeight: monitor.Height,
		Reasoning:       "Default scaling",
	}
}

// ScalingRecommendation contains scaling recommendations for a monitor
type ScalingRecommendation struct {
	MonitorScale    float64
	FontScale       float64
	EffectiveWidth  int
	EffectiveHeight int
	Reasoning       string
}

// ConfigManager handles configuration file management
type ConfigManager struct {
	isDemoMode bool
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(demoMode bool) *ConfigManager {
	return &ConfigManager{
		isDemoMode: demoMode,
	}
}

// ApplyMonitorScale applies monitor scaling to Hyprland configuration
func (cm *ConfigManager) ApplyMonitorScale(monitor Monitor, scale float64) error {
	if cm.isDemoMode {
		// In demo mode, just log what would be done
		fmt.Printf("Demo: Would apply monitor scale %.2fx to %s\n", scale, monitor.Name)
		return nil
	}

	// Check if we're in a Hyprland environment
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		return fmt.Errorf("not running in Hyprland environment")
	}

	// Validate scale - Hyprland prefers clean divisors
	validatedScale := cm.validateHyprlandScale(scale)
	if validatedScale != scale {
		fmt.Printf("Adjusted scale from %.3f to %.3f for Hyprland compatibility\n", scale, validatedScale)
	}

	// Apply scale using hyprctl
	cmd := exec.Command("hyprctl", "keyword", "monitor",
		fmt.Sprintf("%s,preferred,auto,%.5f", monitor.Name, validatedScale))

	// Capture both stdout and stderr to handle errors properly
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Don't show the raw hyprctl error to user, return a clean error
		return fmt.Errorf("failed to apply scaling: Hyprland rejected scale %.3f", validatedScale)
	}

	// Check if output contains any warnings/errors
	outputStr := string(output)
	if strings.Contains(outputStr, "invalid scale") || strings.Contains(outputStr, "failed to find clean divisor") {
		return fmt.Errorf("hyprland rejected scale %.3f - try a different scaling value", validatedScale)
	}

	return nil
}

// validateHyprlandScale ensures the scale value is compatible with Hyprland
func (cm *ConfigManager) validateHyprlandScale(scale float64) float64 {
	// Hyprland prefers these clean divisor scales
	hyprlandScales := []float64{
		1.0,     // 1x
		1.25,    // 5/4
		1.33333, // 4/3
		1.5,     // 3/2
		1.66667, // 5/3
		1.75,    // 7/4
		2.0,     // 2x
		2.25,    // 9/4
		2.5,     // 5/2
		3.0,     // 3x
	}

	// Find the closest valid scale
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

// ApplyGTKScale applies GTK scaling (integer only, as GTK3 doesn't support fractional)
func (cm *ConfigManager) ApplyGTKScale(scale int) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply GTK scale %dx system-wide\n", scale)
		return nil
	}

	// Set GDK_SCALE environment variable
	// This requires logout/login to take full effect
	os.Setenv("GDK_SCALE", fmt.Sprintf("%d", scale))

	// TODO: Could write to ~/.profile or similar for persistence
	return nil
}

// ApplyFontDPI applies font DPI scaling via Xft.dpi
func (cm *ConfigManager) ApplyFontDPI(dpi int) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would set Xft.dpi to %d in ~/.Xresources\n", dpi)
		return nil
	}

	// TODO: Actually implement .Xresources updating
	// For now, just set environment variable for immediate effect
	os.Setenv("XFT_DPI", fmt.Sprintf("%d", dpi))

	return nil
}

// ApplyCompleteScalingOption applies a complete scaling configuration
func (cm *ConfigManager) ApplyCompleteScalingOption(monitor Monitor, option ScalingOption) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply complete scaling option '%s':\n", option.DisplayName)
		fmt.Printf("  - Monitor scale: %.2fx\n", option.MonitorScale)
		fmt.Printf("  - GTK scale: %dx\n", option.GTKScale)
		fmt.Printf("  - Font DPI: %d\n", option.FontDPI)
		fmt.Printf("  - Additional font scale: %.2fx\n", option.FontScale)
		return nil
	}

	// Apply monitor scaling
	if err := cm.ApplyMonitorScale(monitor, option.MonitorScale); err != nil {
		return fmt.Errorf("failed to apply monitor scale: %w", err)
	}

	// Apply GTK scaling
	if err := cm.ApplyGTKScale(option.GTKScale); err != nil {
		return fmt.Errorf("failed to apply GTK scale: %w", err)
	}

	// Apply font DPI
	if err := cm.ApplyFontDPI(option.FontDPI); err != nil {
		return fmt.Errorf("failed to apply font DPI: %w", err)
	}

	return nil
}

// GetScalingExplanations returns explanations for different scaling types
func (cm *ConfigManager) GetScalingExplanations() map[string]string {
	return map[string]string{
		"monitor": "Changes the size of all UI elements rendered by the compositor. Applied immediately, affects everything.",
		"gtk":     "Scales GTK applications (most Linux apps). Requires logout/login. Only supports integer values (1x, 2x, 3x).",
		"font":    "Changes text size system-wide via DPI. Requires reload of applications. Fine-grained control over text rendering.",
	}
}
