package main

import (
	"fmt"
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
			}
		}

		// Resolution line: "2880x1920@120.000Hz at 0x0"
		if strings.Contains(line, "x") && strings.Contains(line, "@") {
			re := regexp.MustCompile(`(\d+)x(\d+)@([\d.]+)Hz at (\d+)x(\d+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 5 {
				currentMonitor.Width, _ = strconv.Atoi(matches[1])
				currentMonitor.Height, _ = strconv.Atoi(matches[2])
				currentMonitor.RefreshRate, _ = strconv.ParseFloat(matches[3], 64)
				currentMonitor.Position.X, _ = strconv.Atoi(matches[4])
				currentMonitor.Position.Y, _ = strconv.Atoi(matches[5])
				currentMonitor.IsActive = true
			}
		}

		// Scale line: "scale: 2.00"
		if strings.Contains(line, "scale:") {
			re := regexp.MustCompile(`scale:\s*([\d.]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentMonitor.Scale, _ = strconv.ParseFloat(matches[1], 64)
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

// GetRecommendedScale returns recommended scaling based on monitor resolution
func (sm *ScalingManager) GetRecommendedScale(monitor Monitor) ScalingRecommendation {
	// Calculate pixel density (rough estimation)
	pixelCount := monitor.Width * monitor.Height

	var monitorScale float64
	var fontScale float64
	var reasoning string

	switch {
	case pixelCount >= 8294400: // 4K+ (3840x2160)
		monitorScale = 2.0
		fontScale = 0.8
		reasoning = "4K+ display: 2x scaling recommended for comfortable viewing"

	case pixelCount >= 3686400: // 1440p (2560x1440)
		monitorScale = 1.0
		fontScale = 0.9
		reasoning = "1440p display: 1x scaling with slightly larger fonts"

	case pixelCount >= 2073600: // 1080p (1920x1080)
		monitorScale = 1.0
		fontScale = 0.8
		reasoning = "1080p display: 1x scaling with standard fonts"

	default:
		monitorScale = 1.0
		fontScale = 1.0
		reasoning = "Standard scaling for this resolution"
	}

	return ScalingRecommendation{
		MonitorScale:    monitorScale,
		FontScale:       fontScale,
		EffectiveWidth:  int(float64(monitor.Width) / monitorScale),
		EffectiveHeight: int(float64(monitor.Height) / monitorScale),
		Reasoning:       reasoning,
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
		fmt.Printf("Demo: Would apply scale %.1fx to monitor %s\n", scale, monitor.Name)
		return nil
	}

	// Check if we're in a Hyprland environment
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		return fmt.Errorf("not running in Hyprland environment")
	}

	// Apply scale using hyprctl
	cmd := exec.Command("hyprctl", "keyword", "monitor",
		fmt.Sprintf("%s,preferred,auto,%.1f", monitor.Name, scale))

	return cmd.Run()
}

// ApplyFontScale applies font scaling across the system
func (cm *ConfigManager) ApplyFontScale(scale float64) error {
	if cm.isDemoMode {
		fmt.Printf("Demo: Would apply font scale %.1fx across the system\n", scale)
		return nil
	}

	// This would apply font scaling to various configuration files
	// For now, just return success in demo mode
	return nil
}
