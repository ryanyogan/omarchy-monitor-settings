package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Terminal-adaptive color scheme using ANSI colors
// Automatically adapts to any terminal theme: Tokyo Night, Catppuccin, GitHub Light, etc.
var (
	// Background colors - use terminal defaults and ANSI colors
	colorBackground = lipgloss.Color("")  // Terminal default background
	colorSurface    = lipgloss.Color("0") // ANSI black (adapts to theme)
	colorFloat      = lipgloss.Color("8") // ANSI bright black/gray

	// Foreground colors - use terminal defaults and ANSI colors
	colorForeground = lipgloss.Color("")  // Terminal default foreground
	colorComment    = lipgloss.Color("8") // ANSI bright black (dim)
	colorSubtle     = lipgloss.Color("7") // ANSI white (adapts to theme)

	// Accent colors - standard ANSI that adapt to any theme
	colorBlue    = lipgloss.Color("4") // ANSI blue
	colorCyan    = lipgloss.Color("6") // ANSI cyan
	colorGreen   = lipgloss.Color("2") // ANSI green
	colorYellow  = lipgloss.Color("3") // ANSI yellow
	colorRed     = lipgloss.Color("1") // ANSI red
	colorMagenta = lipgloss.Color("5") // ANSI magenta
)

// getTerminalThemeInfo returns information about the current terminal theme
func getTerminalThemeInfo() string {
	termOutput := termenv.NewOutput(os.Stdout)
	profile := termOutput.Profile
	isDark := termOutput.HasDarkBackground()

	var profileName string
	switch profile {
	case termenv.TrueColor:
		profileName = "TrueColor (24-bit)"
	case termenv.ANSI256:
		profileName = "256 Color"
	case termenv.ANSI:
		profileName = "16 Color"
	default:
		profileName = "Basic"
	}

	theme := "Unknown"
	if isDark {
		theme = "Dark"
	} else {
		theme = "Light"
	}

	return fmt.Sprintf("Terminal Adaptive (%s, %s)", profileName, theme)
}

// getValidHyprlandScales returns the scales that Hyprland accepts without errors
func getValidHyprlandScales() []float64 {
	return []float64{1.0, 1.25, 1.33333, 1.5, 1.66667, 1.75, 2.0, 2.25, 2.5, 3.0}
}

// findNextValidScale finds the next valid scale in the direction specified
func findNextValidScale(current float64, up bool) float64 {
	validScales := getValidHyprlandScales()

	// Find current position
	currentIndex := -1
	for i, scale := range validScales {
		if math.Abs(scale-current) < 0.001 {
			currentIndex = i
			break
		}
	}

	// If not found, find closest
	if currentIndex == -1 {
		minDiff := math.Abs(validScales[0] - current)
		currentIndex = 0
		for i, scale := range validScales {
			diff := math.Abs(scale - current)
			if diff < minDiff {
				minDiff = diff
				currentIndex = i
			}
		}
	}

	// Move in the specified direction
	if up {
		if currentIndex < len(validScales)-1 {
			return validScales[currentIndex+1]
		}
		return validScales[len(validScales)-1] // Stay at max
	} else {
		if currentIndex > 0 {
			return validScales[currentIndex-1]
		}
		return validScales[0] // Stay at min
	}
}

// AppMode represents different screens/modes in the TUI
type AppMode int

const (
	ModeDashboard AppMode = iota
	ModeMonitorSelection
	ModeScalingOptions
	ModeManualScaling
	ModeSettings
	ModeHelp
	ModeConfirmation
)

// ConfirmationAction represents what action is pending confirmation
type ConfirmationAction int

const (
	ConfirmNone ConfirmationAction = iota
	ConfirmSmartScaling
	ConfirmManualScaling
)

// Monitor represents a detected monitor
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

// Model is the main bubbletea model
type Model struct {
	// App state
	mode   AppMode
	width  int
	height int
	ready  bool

	// Services (injected dependencies)
	services *AppServices

	// Monitor data
	monitors        []Monitor
	selectedMonitor int
	isDemoMode      bool

	// Scaling data
	scalingOptions     []ScalingOption
	selectedScalingOpt int

	// Manual scaling controls
	manualMonitorScale    float64
	manualGTKScale        int
	manualFontDPI         int
	selectedManualControl int // 0=Monitor Scale, 1=GTK Scale, 2=Font DPI

	// Confirmation state
	confirmationAction   ConfirmationAction
	pendingScalingOption ScalingOption
	pendingMonitor       Monitor

	// UI state
	selectedOption int
	menuItems      []string

	// Styles
	headerStyle     lipgloss.Style
	footerStyle     lipgloss.Style
	titleStyle      lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	helpStyle       lipgloss.Style
	errorStyle      lipgloss.Style
	successStyle    lipgloss.Style
}

// NewModel creates and initializes a new Model (legacy function for backward compatibility)
func NewModel() Model {
	// Create default services for backward compatibility
	config := &AppConfig{
		NoHyprlandCheck: noHyprlandCheck,
		DebugMode:       debugMode,
		ForceLiveMode:   forceLiveMode,
		IsTestMode:      false,
	}
	services := NewAppServices(config)
	return NewModelWithServices(services)
}

// NewModelWithServices creates and initializes a new Model with injected services
func NewModelWithServices(services *AppServices) Model {
	m := Model{
		mode:           ModeDashboard,
		selectedOption: 0,
		menuItems: []string{
			"Dashboard",
			"Monitor Selection",
			"Smart Scaling",
			"Manual Scaling",
			"Settings",
			"Help",
			"Exit",
		},
		isDemoMode: true, // Default to demo mode for testing
		services:   services,

		// Initialize manual scaling defaults
		manualMonitorScale:    1.0,
		manualGTKScale:        1,
		manualFontDPI:         96,
		selectedManualControl: 0,
	}

	// Initialize styles with Tokyo Night theme
	m.initStyles()

	// Load monitors using detection or demo data
	m.loadMonitors()

	// Load intelligent scaling options for the first monitor
	if len(m.monitors) > 0 {
		m.scalingOptions = services.ScalingManager.GetIntelligentScalingOptions(m.monitors[0])
	}

	return m
}

func (m *Model) initStyles() {
	// Header style - cleaner btop-like header
	m.headerStyle = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorForeground).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue)

	// Footer style
	m.footerStyle = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorComment).
		Padding(0, 1)

	// Title style
	m.titleStyle = lipgloss.NewStyle().
		Foreground(colorBlue).
		Bold(true).
		Underline(true)

		// Selected menu item style - no background, just text styling
	m.selectedStyle = lipgloss.NewStyle().
		Foreground(colorBlue).
		Bold(true)

	// Unselected menu item style
	m.unselectedStyle = lipgloss.NewStyle().
		Foreground(colorForeground)

	// Help text style
	m.helpStyle = lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true)

	// Error style
	m.errorStyle = lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)

	// Success style
	m.successStyle = lipgloss.NewStyle().
		Foreground(colorGreen).
		Bold(true)
}

func (m *Model) loadMonitors() {
	// Use injected monitor detector
	monitors, err := m.services.MonitorDetector.DetectMonitors()

	if m.services.Config.DebugMode {
		fmt.Printf("DEBUG: DetectMonitors returned %d monitors, error: %v\n", len(monitors), err)
	}

	if err != nil {
		// Fallback to demo monitors
		if m.services.Config.DebugMode {
			fmt.Printf("DEBUG: Setting demo mode due to detection error\n")
		}
		m.isDemoMode = true
		// Get fallback monitors from the detector
		if detector, ok := m.services.MonitorDetector.(*MonitorDetector); ok {
			monitors, _ = detector.getFallbackMonitors()
		}
	} else {
		if m.services.Config.DebugMode {
			fmt.Printf("DEBUG: Setting live mode - detected real monitors\n")
		}
		m.isDemoMode = false
	}

	// Override demo mode if force-live flag is set
	if m.services.Config.ForceLiveMode {
		if m.services.Config.DebugMode {
			fmt.Printf("DEBUG: Force-live mode enabled, overriding demo mode\n")
		}
		m.isDemoMode = false
	}

	if m.services.Config.DebugMode {
		fmt.Printf("DEBUG: Final state - isDemoMode: %v, monitor count: %d\n", m.isDemoMode, len(monitors))
	}

	m.monitors = monitors
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	default:
		return m, nil
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		switch m.mode {
		case ModeDashboard:
			if m.selectedOption > 0 {
				m.selectedOption--
			}
		case ModeMonitorSelection:
			if m.selectedMonitor > 0 {
				m.selectedMonitor--
				// Update scaling options when monitor changes
				if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
					m.scalingOptions = m.services.ScalingManager.GetIntelligentScalingOptions(m.monitors[m.selectedMonitor])
					m.selectedScalingOpt = 0 // Reset to first option
				}
			}
		case ModeScalingOptions:
			if m.selectedScalingOpt > 0 {
				m.selectedScalingOpt--
			}
		case ModeManualScaling:
			// Navigate to previous control
			if m.selectedManualControl > 0 {
				m.selectedManualControl--
			}
		}

	case "down", "j":
		switch m.mode {
		case ModeDashboard:
			if m.selectedOption < len(m.menuItems)-1 {
				m.selectedOption++
			}
		case ModeMonitorSelection:
			if m.selectedMonitor < len(m.monitors)-1 {
				m.selectedMonitor++
				// Update scaling options when monitor changes
				if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
					m.scalingOptions = m.services.ScalingManager.GetIntelligentScalingOptions(m.monitors[m.selectedMonitor])
					m.selectedScalingOpt = 0 // Reset to first option
				}
			}
		case ModeScalingOptions:
			if m.selectedScalingOpt < len(m.scalingOptions)-1 {
				m.selectedScalingOpt++
			}
		case ModeManualScaling:
			// Navigate to next control
			if m.selectedManualControl < 2 {
				m.selectedManualControl++
			}
		}

	case "left":
		if m.mode == ModeManualScaling {
			// Decrease the selected manual control value
			switch m.selectedManualControl {
			case 0: // Monitor Scale - use Hyprland-compatible scales only
				m.manualMonitorScale = findNextValidScale(m.manualMonitorScale, false)
			case 1: // GTK Scale
				if m.manualGTKScale > 1 {
					m.manualGTKScale--
				}
			case 2: // Font DPI
				if m.manualFontDPI > 72 {
					m.manualFontDPI -= 12
					if m.manualFontDPI < 72 {
						m.manualFontDPI = 72
					}
				}
			}
		}

	case "right":
		if m.mode == ModeManualScaling {
			// Increase the selected manual control value
			switch m.selectedManualControl {
			case 0: // Monitor Scale - use Hyprland-compatible scales only
				m.manualMonitorScale = findNextValidScale(m.manualMonitorScale, true)
			case 1: // GTK Scale
				if m.manualGTKScale < 3 {
					m.manualGTKScale++
				}
			case 2: // Font DPI
				if m.manualFontDPI < 288 {
					m.manualFontDPI += 12
					if m.manualFontDPI > 288 {
						m.manualFontDPI = 288
					}
				}
			}
		}

	case "enter", " ":
		if m.mode == ModeMonitorSelection {
			// Monitor selected - return to dashboard
			m.mode = ModeDashboard
			m.selectedOption = 0
			return m, nil
		} else if m.mode == ModeScalingOptions && len(m.scalingOptions) > 0 && m.selectedScalingOpt < len(m.scalingOptions) {
			// Show confirmation for smart scaling
			selectedOption := m.scalingOptions[m.selectedScalingOpt]
			if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
				m.confirmationAction = ConfirmSmartScaling
				m.pendingScalingOption = selectedOption
				m.pendingMonitor = m.monitors[m.selectedMonitor]
				m.mode = ModeConfirmation
			}
			return m, nil
		} else if m.mode == ModeManualScaling {
			// Show confirmation for manual scaling
			if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
				m.confirmationAction = ConfirmManualScaling
				m.pendingMonitor = m.monitors[m.selectedMonitor]
				// Create a temporary scaling option for manual settings
				m.pendingScalingOption = ScalingOption{
					MonitorScale: m.manualMonitorScale,
					GTKScale:     m.manualGTKScale,
					FontDPI:      m.manualFontDPI,
					DisplayName:  "Manual Settings",
					Description:  "Custom scaling values",
				}
				m.mode = ModeConfirmation
			}
			return m, nil
		} else if m.mode == ModeConfirmation {
			// Apply the confirmed scaling using injected config manager
			switch m.confirmationAction {
			case ConfirmSmartScaling:
				m.services.ConfigManager.ApplyCompleteScalingOption(m.pendingMonitor, m.pendingScalingOption)
			case ConfirmManualScaling:
				m.services.ConfigManager.ApplyMonitorScale(m.pendingMonitor, m.manualMonitorScale)
				m.services.ConfigManager.ApplyGTKScale(m.manualGTKScale)
				m.services.ConfigManager.ApplyFontDPI(m.manualFontDPI)
			}
			// Reset confirmation state and return to dashboard
			m.confirmationAction = ConfirmNone
			m.mode = ModeDashboard
			m.selectedOption = 0
			return m, nil
		}
		return m.handleSelection()

	case "m":
		if m.mode == ModeScalingOptions {
			m.mode = ModeManualScaling
			return m, nil
		}

	case "h", "?":
		m.mode = ModeHelp

	case "esc":
		switch m.mode {
		case ModeManualScaling:
			m.mode = ModeDashboard
			m.selectedOption = 0 // Reset to first option
		case ModeConfirmation:
			// Cancel confirmation and return to previous mode
			switch m.confirmationAction {
			case ConfirmSmartScaling:
				m.mode = ModeScalingOptions
			case ConfirmManualScaling:
				m.mode = ModeManualScaling
			default:
				m.mode = ModeDashboard
			}
			m.confirmationAction = ConfirmNone
		default:
			m.mode = ModeDashboard
			// Reset selection when returning to dashboard
			m.selectedOption = 0
		}
	}

	return m, nil
}

func (m Model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.selectedOption {
	case 0: // Dashboard
		m.mode = ModeDashboard
	case 1: // Monitor Selection
		m.mode = ModeMonitorSelection
	case 2: // Smart Scaling
		m.mode = ModeScalingOptions
	case 3: // Manual Scaling
		m.mode = ModeManualScaling
	case 4: // Settings
		m.mode = ModeSettings
	case 5: // Help
		m.mode = ModeHelp
	case 6: // Exit
		return m, tea.Quit
	}

	return m, nil
}

// View renders the award-winning beautiful TUI
func (m Model) View() string {
	// Ensure minimum dimensions
	if m.width < 80 || m.height < 20 {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorRed).
			Render("Terminal too small\nPlease resize to at least 80x20")
	}

	if !m.ready {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorBlue).
			Render("Initializing stunning TUI...")
	}

	// Calculate precise dimensions for clean layout with simple header
	headerHeight := 5                                           // Header with border, padding, and margin
	footerHeight := 2                                           // Footer spacing
	contentHeight := m.height - headerHeight - footerHeight - 2 // -2 for margins

	// Ensure content doesn't overflow
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Build the view based on current mode
	var content string

	switch m.mode {
	case ModeDashboard:
		content = m.renderDashboard(contentHeight)
	case ModeMonitorSelection:
		content = m.renderMonitorSelection(contentHeight)
	case ModeScalingOptions:
		content = m.renderScalingOptions(contentHeight)
	case ModeManualScaling:
		content = m.renderManualScaling(contentHeight)
	case ModeSettings:
		content = m.renderSettings(contentHeight)
	case ModeHelp:
		content = m.renderHelp(contentHeight)
	case ModeConfirmation:
		content = m.renderConfirmation(contentHeight)
	default:
		content = m.renderDashboard(contentHeight)
	}

	// Create clean layout with simple header
	header := m.renderHeader()
	footer := m.renderFooter()

	// Professional content container with perfect proportions
	styledContent := lipgloss.NewStyle().
		Width(m.width-4). // Leave margin on sides
		Height(contentHeight).
		Margin(1, 2). // Perfect margins
		Render(content)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		styledContent,
		footer,
	)
}

func (m Model) renderHeader() string {
	// Ultra simple header to test text visibility
	return lipgloss.NewStyle().
		Width(m.width).
		Background(colorBackground).
		Foreground(colorBlue).
		Bold(true).
		Align(lipgloss.Center).
		Render("Display Settings")
}

func (m Model) renderFooter() string {
	// Award-winning footer with elegant key hints
	keyStyle := lipgloss.NewStyle().
		Background(colorBackground).
		Foreground(colorBlue).
		Bold(true).
		Padding(0, 1).
		Margin(0, 1)

	textStyle := lipgloss.NewStyle().
		Foreground(colorSubtle)

	controls := []string{
		keyStyle.Foreground(colorGreen).Render("↑↓") + textStyle.Render("navigate"),
		keyStyle.Foreground(colorBlue).Render("⏎") + textStyle.Render("select"),
		keyStyle.Foreground(colorYellow).Render("h") + textStyle.Render("help"),
		keyStyle.Foreground(colorMagenta).Render("esc") + textStyle.Render("back"),
		keyStyle.Foreground(colorRed).Render("q") + textStyle.Render("quit"),
	}

	helpText := strings.Join(controls, "  ")

	// Calculate the total width of the dashboard content (both panels + gap)
	availableWidth := m.width - 8                // Account for margins and borders
	leftWidth := availableWidth * 2 / 5          // 40% for menu
	rightWidth := availableWidth - leftWidth - 4 // Remaining for monitors

	// Ensure minimum widths
	if leftWidth < 25 {
		leftWidth = 25
	}
	if rightWidth < 30 {
		rightWidth = 30
	}

	// Total footer width = left panel + gap + right panel
	// Need to account for the fact that panels are joined horizontally with a 2-space gap
	// Plus account for the content margin (2 spaces on each side)
	totalFooterWidth := leftWidth + 2 + rightWidth + 2 // +2 for content margin (1 on each side)

	return lipgloss.NewStyle().
		Width(totalFooterWidth).
		Background(colorBackground).
		Align(lipgloss.Center).
		Padding(1, 2).
		Margin(0, 2). // Reduced margin to align with panel borders, not internal text
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorComment).
		Render(helpText)
}

func (m Model) renderDashboard(contentHeight int) string {
	// Award-winning responsive layout with perfect proportions
	availableWidth := m.width - 8                // Account for margins and borders
	leftWidth := availableWidth * 2 / 5          // 40% for menu
	rightWidth := availableWidth - leftWidth - 4 // Remaining for monitors

	// Ensure minimum widths
	if leftWidth < 25 {
		leftWidth = 25
	}
	if rightWidth < 30 {
		rightWidth = 30
	}

	var leftPanel []string
	var rightPanel []string

	// Elegant menu panel
	leftPanel = append(leftPanel,
		lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true).
			Render("Navigation"),
	)
	leftPanel = append(leftPanel, "")

	menuColors := []lipgloss.Color{colorBlue, colorCyan, colorGreen, colorYellow, colorRed, colorMagenta}

	for i, item := range m.menuItems {
		color := menuColors[i%len(menuColors)]
		var line string

		if i == m.selectedOption {
			// Beautiful selection indicator
			selector := lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render("▶")
			text := lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render(item)
			line = fmt.Sprintf("%s %s", selector, text)
		} else {
			text := lipgloss.NewStyle().
				Foreground(colorSubtle).
				Render(item)
			line = fmt.Sprintf("  %s", text)
		}
		leftPanel = append(leftPanel, line)
		leftPanel = append(leftPanel, "") // Breathing room
	}

	// Elegant monitor panel
	rightPanel = append(rightPanel,
		lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true).
			Render("Display Overview"),
	)
	rightPanel = append(rightPanel, "")

	monitorColors := []lipgloss.Color{colorGreen, colorBlue, colorYellow, colorMagenta}

	for i, monitor := range m.monitors {
		if i >= len(monitorColors) {
			break // Prevent overflow
		}

		color := monitorColors[i]

		// Elegant status indicators
		var statusIcon string
		var statusStyle lipgloss.Style
		if monitor.IsActive {
			if monitor.IsPrimary {
				statusIcon = "●"
				statusStyle = lipgloss.NewStyle().Foreground(colorGreen)
			} else {
				statusIcon = "○"
				statusStyle = lipgloss.NewStyle().Foreground(colorBlue)
			}
		} else {
			statusIcon = "◦"
			statusStyle = lipgloss.NewStyle().Foreground(colorComment)
		}

		// Beautiful monitor card
		header := fmt.Sprintf("%s %s",
			statusStyle.Render(statusIcon),
			lipgloss.NewStyle().Foreground(color).Bold(true).Render(monitor.Name),
		)

		// Add indicator if this is the currently selected monitor for changes
		if i == m.selectedMonitor {
			selectedIndicator := lipgloss.NewStyle().
				Foreground(colorYellow).
				Bold(true).
				Render(" 👆 CURRENT")
			header = header + selectedIndicator
		}

		details := []string{
			lipgloss.NewStyle().Foreground(colorSubtle).Render(fmt.Sprintf("  %s %s", monitor.Make, monitor.Model)),
			lipgloss.NewStyle().Foreground(colorComment).Render(fmt.Sprintf("  %dx%d @ %.0fHz", monitor.Width, monitor.Height, monitor.RefreshRate)),
			lipgloss.NewStyle().Foreground(colorComment).Render(fmt.Sprintf("  Scale: %.1fx", monitor.Scale)),
		}

		// Add extra info if this is the selected monitor
		if i == m.selectedMonitor {
			details = append(details, lipgloss.NewStyle().
				Foreground(colorYellow).
				Italic(true).
				Render("  → Scaling changes will apply here"))
		}

		rightPanel = append(rightPanel, header)
		rightPanel = append(rightPanel, details...)
		rightPanel = append(rightPanel, "") // Card spacing
	}

	if m.isDemoMode {
		demoNotice := lipgloss.NewStyle().
			Foreground(colorYellow).
			Italic(true).
			Render("�� Demo Mode Active")
		rightPanel = append(rightPanel, demoNotice)
	}

	// Award-winning panel styling
	leftContent := lipgloss.NewStyle().
		Width(leftWidth).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Render(strings.Join(leftPanel, "\n"))

	rightContent := lipgloss.NewStyle().
		Width(rightWidth).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorCyan).
		Render(strings.Join(rightPanel, "\n"))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftContent,
		lipgloss.NewStyle().Width(2).Render(""), // Perfect spacing
		rightContent,
	)
}

func (m Model) renderMonitorSelection(contentHeight int) string {
	var content []string

	// Award-winning monitor selection screen
	title := lipgloss.NewStyle().
		Foreground(colorBlue).
		Bold(true).
		Render("Monitor Selection")

	subtitle := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Render("Choose a display to configure")

	content = append(content, title)
	content = append(content, subtitle)
	content = append(content, "")

	monitorColors := []lipgloss.Color{colorGreen, colorBlue, colorYellow, colorMagenta}

	for i, monitor := range m.monitors {
		if i >= len(monitorColors) {
			break
		}

		color := monitorColors[i]

		// Status styling
		var statusText string
		var statusStyle lipgloss.Style
		if monitor.IsActive {
			if monitor.IsPrimary {
				statusText = "PRIMARY"
				statusStyle = lipgloss.NewStyle().Background(colorGreen).Foreground(colorBackground).Bold(true).Padding(0, 1)
			} else {
				statusText = "ACTIVE"
				statusStyle = lipgloss.NewStyle().Background(colorBlue).Foreground(colorBackground).Bold(true).Padding(0, 1)
			}
		} else {
			statusText = "INACTIVE"
			statusStyle = lipgloss.NewStyle().Background(colorComment).Foreground(colorBackground).Bold(true).Padding(0, 1)
		}

		// Beautiful monitor card
		nameStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
		detailStyle := lipgloss.NewStyle().Foreground(colorSubtle)

		var card string
		if i == m.selectedMonitor {
			selector := lipgloss.NewStyle().Foreground(color).Bold(true).Render("▶ ")
			card = fmt.Sprintf("%s%s %s\n  %s\n  %s",
				selector,
				nameStyle.Render(monitor.Name),
				statusStyle.Render(statusText),
				detailStyle.Render(fmt.Sprintf("%s %s", monitor.Make, monitor.Model)),
				detailStyle.Render(fmt.Sprintf("%dx%d @ %.0fHz", monitor.Width, monitor.Height, monitor.RefreshRate)),
			)
		} else {
			card = fmt.Sprintf("  %s %s\n    %s\n    %s",
				nameStyle.Render(monitor.Name),
				statusStyle.Render(statusText),
				detailStyle.Render(fmt.Sprintf("%s %s", monitor.Make, monitor.Model)),
				detailStyle.Render(fmt.Sprintf("%dx%d @ %.0fHz", monitor.Width, monitor.Height, monitor.RefreshRate)),
			)
		}

		content = append(content, card)
		content = append(content, "")
	}

	// Instructions
	instructions := []string{
		lipgloss.NewStyle().Foreground(colorYellow).Render("⏎") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" Select monitor and return to dashboard"),
		lipgloss.NewStyle().Foreground(colorMagenta).Render("esc") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" Return to main menu"),
	}

	// Add helpful note
	note := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true).
		Render("💡 Selected monitor will be marked as CURRENT on the dashboard")

	content = append(content, "")
	content = append(content, strings.Join(instructions, "  "))
	content = append(content, "")
	content = append(content, note)

	// Beautiful container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderScalingOptions(contentHeight int) string {
	var content []string

	// Award-winning smart scaling screen
	title := lipgloss.NewStyle().
		Foreground(colorGreen).
		Bold(true).
		Render("🧠 Smart Scaling Recommendations")

	content = append(content, title)
	content = append(content, "")

	if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
		selectedMonitor := m.monitors[m.selectedMonitor]

		// Monitor info card - simplified
		monitorInfo := fmt.Sprintf("%s %s %s - %dx%d@%.0fHz",
			selectedMonitor.Name, selectedMonitor.Make, selectedMonitor.Model,
			selectedMonitor.Width, selectedMonitor.Height, selectedMonitor.RefreshRate)

		monitorCard := lipgloss.NewStyle().
			Background(colorSurface).
			Foreground(colorForeground).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorYellow).
			Render(monitorInfo)

		content = append(content, monitorCard)
		content = append(content, "")

		// Smart recommendations - simplified layout
		recTitle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("🎯 Available Options")
		content = append(content, recTitle)
		content = append(content, "")

		// Display scaling options in a simpler format
		for i, option := range m.scalingOptions {
			var line string

			// Selection indicator and recommended badge
			if i == m.selectedScalingOpt {
				line = lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("▶ ")
			} else {
				line = "  "
			}

			// Option name with recommended badge
			optionName := option.DisplayName
			if option.IsRecommended {
				optionName += " [RECOMMENDED]"
			}

			if i == m.selectedScalingOpt {
				line += lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render(optionName)
			} else {
				line += lipgloss.NewStyle().Foreground(colorForeground).Render(optionName)
			}

			content = append(content, line)

			// Description
			description := fmt.Sprintf("    %s", option.Description)
			content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(description))

			// Technical details
			details := fmt.Sprintf("    Monitor: %.1fx • GTK: %dx • Font DPI: %d • Result: %dx%d",
				option.MonitorScale, option.GTKScale, option.FontDPI,
				option.EffectiveWidth, option.EffectiveHeight)
			content = append(content, lipgloss.NewStyle().Foreground(colorComment).Render(details))

			// Reasoning
			reasoning := fmt.Sprintf("    💡 %s", option.Reasoning)
			content = append(content, lipgloss.NewStyle().Foreground(colorComment).Italic(true).Render(reasoning))

			content = append(content, "") // Space between options
		}

		// What each setting does - simplified
		content = append(content, "")
		explainTitle := lipgloss.NewStyle().Foreground(colorMagenta).Bold(true).Render("📚 What Each Setting Does")
		content = append(content, explainTitle)

		explainItems := []string{
			"Monitor Scale: Changes compositor-level scaling (immediate effect)",
			"GTK Scale: Scales GTK applications (requires logout/login)",
			"Font DPI: Fine-grained text scaling (affects most apps)",
		}

		for _, item := range explainItems {
			content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render("  "+item))
		}

		if m.isDemoMode {
			content = append(content, "")
			demoNotice := lipgloss.NewStyle().
				Foreground(colorYellow).
				Render("📱 Demo Mode: Changes will be simulated")
			content = append(content, demoNotice)
		}
	}

	// Instructions
	content = append(content, "")
	instructions := []string{
		lipgloss.NewStyle().Foreground(colorGreen).Render("↑↓") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" select"),
		lipgloss.NewStyle().Foreground(colorBlue).Render("⏎") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" apply"),
		lipgloss.NewStyle().Foreground(colorYellow).Render("m") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" manual"),
		lipgloss.NewStyle().Foreground(colorMagenta).Render("esc") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" back"),
	}

	content = append(content, strings.Join(instructions, "  "))

	// Simple container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorGreen).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderManualScaling(contentHeight int) string {
	var content []string

	// Award-winning manual scaling screen
	title := lipgloss.NewStyle().
		Foreground(colorMagenta).
		Bold(true).
		Render("🔧 Manual Scaling Controls")

	content = append(content, title)
	content = append(content, "")

	// Check bounds to prevent crashes
	if len(m.monitors) == 0 {
		content = append(content, lipgloss.NewStyle().
			Foreground(colorRed).
			Render("No monitors detected. Please go back and check monitor selection."))
		return lipgloss.NewStyle().
			Width(m.width - 8).
			Height(contentHeight - 2).
			Padding(2).
			Background(colorBackground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMagenta).
			Render(strings.Join(content, "\n"))
	}

	if m.selectedMonitor >= len(m.monitors) {
		content = append(content, lipgloss.NewStyle().
			Foreground(colorRed).
			Render("Invalid monitor selection. Please go back and select a monitor."))
		return lipgloss.NewStyle().
			Width(m.width - 8).
			Height(contentHeight - 2).
			Padding(2).
			Background(colorBackground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMagenta).
			Render(strings.Join(content, "\n"))
	}

	selectedMonitor := m.monitors[m.selectedMonitor]

	// Monitor info card - simplified
	monitorInfo := fmt.Sprintf("%s %s %s - %dx%d@%.0fHz (Current: %.2fx)",
		selectedMonitor.Name, selectedMonitor.Make, selectedMonitor.Model,
		selectedMonitor.Width, selectedMonitor.Height, selectedMonitor.RefreshRate, selectedMonitor.Scale)

	monitorCard := lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorForeground).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorYellow).
		Render(monitorInfo)

	content = append(content, monitorCard)
	content = append(content, "")

	// Manual controls - simplified layout with selection indicators
	controlsTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("⚙️ Scaling Controls")
	content = append(content, controlsTitle)
	content = append(content, "")

	// Monitor Scale Control
	monitorScaleStyle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true)
	monitorScaleValueStyle := lipgloss.NewStyle().Foreground(colorGreen)
	monitorScaleDescStyle := lipgloss.NewStyle().Foreground(colorSubtle)

	if m.selectedManualControl == 0 {
		monitorScaleStyle = monitorScaleStyle.Background(colorSurface).Padding(0, 1)
		monitorScaleValueStyle = monitorScaleValueStyle.Background(colorSurface).Padding(0, 1)
		monitorScaleDescStyle = monitorScaleDescStyle.Background(colorSurface).Padding(0, 1)
	}

	monitorScaleLabel := monitorScaleStyle.Render("1. Monitor Scale (Compositor-level)")
	content = append(content, monitorScaleLabel)
	monitorScaleValue := fmt.Sprintf("   Current: %.3fx (Valid: 1.0x, 1.25x, 1.33x, 1.5x, 1.67x, 1.75x, 2.0x, 2.25x, 2.5x, 3.0x)", m.manualMonitorScale)
	content = append(content, monitorScaleValueStyle.Render(monitorScaleValue))
	content = append(content, monitorScaleDescStyle.Render("   Scales everything immediately. Works with all apps."))
	content = append(content, "")

	// GTK Scale Control
	gtkScaleStyle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	gtkScaleValueStyle := lipgloss.NewStyle().Foreground(colorGreen)
	gtkScaleDescStyle := lipgloss.NewStyle().Foreground(colorSubtle)

	if m.selectedManualControl == 1 {
		gtkScaleStyle = gtkScaleStyle.Background(colorSurface).Padding(0, 1)
		gtkScaleValueStyle = gtkScaleValueStyle.Background(colorSurface).Padding(0, 1)
		gtkScaleDescStyle = gtkScaleDescStyle.Background(colorSurface).Padding(0, 1)
	}

	gtkScaleLabel := gtkScaleStyle.Render("2. GTK Scale (Application-level)")
	content = append(content, gtkScaleLabel)
	gtkScaleValue := fmt.Sprintf("   Current: %dx (Range: 1x - 3x, Integer only)", m.manualGTKScale)
	content = append(content, gtkScaleValueStyle.Render(gtkScaleValue))
	content = append(content, gtkScaleDescStyle.Render("   Scales GTK apps (most Linux apps). Requires logout."))
	content = append(content, "")

	// Font DPI Control
	fontDPIStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
	fontDPIValueStyle := lipgloss.NewStyle().Foreground(colorGreen)
	fontDPIDescStyle := lipgloss.NewStyle().Foreground(colorSubtle)

	if m.selectedManualControl == 2 {
		fontDPIStyle = fontDPIStyle.Background(colorSurface).Padding(0, 1)
		fontDPIValueStyle = fontDPIValueStyle.Background(colorSurface).Padding(0, 1)
		fontDPIDescStyle = fontDPIDescStyle.Background(colorSurface).Padding(0, 1)
	}

	fontDPILabel := fontDPIStyle.Render("3. Font DPI (Text rendering)")
	content = append(content, fontDPILabel)
	fontDPIValue := fmt.Sprintf("   Current: %d (Range: 72 - 288, Step: 12)", m.manualFontDPI)
	content = append(content, fontDPIValueStyle.Render(fontDPIValue))
	content = append(content, fontDPIDescStyle.Render("   Fine-grained text scaling. Works with most applications."))
	content = append(content, "")

	// Results preview
	effectiveWidth := int(float64(selectedMonitor.Width) / m.manualMonitorScale)
	effectiveHeight := int(float64(selectedMonitor.Height) / m.manualMonitorScale)
	screenRealEstate := 100.0 / m.manualMonitorScale
	fontMultiplier := float64(m.manualFontDPI) / 96.0

	resultsTitle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("📊 Preview Results")
	content = append(content, resultsTitle)
	content = append(content, fmt.Sprintf("  Effective Resolution: %dx%d", effectiveWidth, effectiveHeight))
	content = append(content, fmt.Sprintf("  Screen Real Estate: %.0f%%", screenRealEstate))
	content = append(content, fmt.Sprintf("  Font DPI Multiplier: %.1fx", fontMultiplier))

	if m.isDemoMode {
		content = append(content, "")
		demoNotice := lipgloss.NewStyle().
			Foreground(colorYellow).
			Render("📱 Demo Mode: Use ⏎ to preview changes")
		content = append(content, demoNotice)
	}

	// Instructions
	content = append(content, "")
	instructions := []string{
		lipgloss.NewStyle().Foreground(colorGreen).Render("↑↓") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" select control"),
		lipgloss.NewStyle().Foreground(colorCyan).Render("←→") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" adjust value"),
		lipgloss.NewStyle().Foreground(colorYellow).Render("⏎") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" apply all"),
		lipgloss.NewStyle().Foreground(colorMagenta).Render("esc") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" back"),
	}

	content = append(content, strings.Join(instructions, "  "))

	// Simple container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMagenta).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderSettings(contentHeight int) string {
	var content []string

	// Award-winning settings screen
	title := lipgloss.NewStyle().
		Foreground(colorMagenta).
		Bold(true).
		Render("⚙️ Application Settings")

	content = append(content, title)
	content = append(content, "")

	// Application info section - clean and consistent
	appTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("📱 Application Info")
	content = append(content, appTitle)
	content = append(content, "")

	appItems := []string{
		fmt.Sprintf("  Version: %s", lipgloss.NewStyle().Foreground(colorGreen).Render("1.0.0")),
		fmt.Sprintf("  Theme: %s", lipgloss.NewStyle().Foreground(colorMagenta).Render(getTerminalThemeInfo())),
		fmt.Sprintf("  Mode: %s", func() string {
			if m.isDemoMode {
				return lipgloss.NewStyle().Foreground(colorYellow).Render("Demo")
			}
			return lipgloss.NewStyle().Foreground(colorGreen).Render("Live")
		}()),
	}

	for _, item := range appItems {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// Detection methods section - clean and consistent
	detectionTitle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("🔍 Detection Methods")
	content = append(content, detectionTitle)
	content = append(content, "")

	detector := NewMonitorDetector()
	methods := []struct {
		name string
		cmd  string
	}{
		{"Hyprctl", "hyprctl"},
		{"wlr-randr", "wlr-randr"},
	}

	for _, method := range methods {
		var status string
		if detector.commandExists(method.cmd) {
			status = lipgloss.NewStyle().Foreground(colorGreen).Render("✓ Available")
		} else {
			status = lipgloss.NewStyle().Foreground(colorRed).Render("✗ Not found")
		}
		item := fmt.Sprintf("  %s: %s", method.name, status)
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// Configuration section - clean and consistent (FIXED)
	configTitle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("⚙️ Configuration")
	content = append(content, configTitle)
	content = append(content, "")

	configItems := []string{
		fmt.Sprintf("  Target: %s", lipgloss.NewStyle().Foreground(colorGreen).Render("Hyprland + Wayland")),
		fmt.Sprintf("  Fallbacks: %s", lipgloss.NewStyle().Foreground(colorBlue).Render("wlr-randr")),
		fmt.Sprintf("  Font Scaling: %s", lipgloss.NewStyle().Foreground(colorMagenta).Render("GTK, Alacritty, Neovim")),
	}

	for _, item := range configItems {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// Footer message
	footer := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true).
		Render("💡 Press Esc to return to the main menu")
	content = append(content, footer)

	// Simple, clean container - consistent with other screens
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMagenta).
		Render(strings.Join(content, "\n"))
}

// renderConfirmation renders the confirmation dialog for scaling operations
func (m Model) renderConfirmation(contentHeight int) string {
	var content []string

	// Title with warning icon
	title := lipgloss.NewStyle().
		Foreground(colorYellow).
		Bold(true).
		Render("⚠️ Confirm Scaling Changes")

	content = append(content, title)
	content = append(content, "")

	// Warning message
	warningStyle := lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)

	warning := warningStyle.Render("⚠️  WARNING: Desktop refresh required!")
	content = append(content, warning)
	content = append(content, "")

	// Explanation
	explanationLines := []string{
		"Applying these changes will:",
		"",
		"🔄 Refresh your desktop environment",
		"❌ Close all applications WITHOUT saving",
		"⚡ Apply new scaling immediately",
		"",
		"Make sure to save your work before proceeding!",
	}

	for _, line := range explanationLines {
		if strings.HasPrefix(line, "🔄") || strings.HasPrefix(line, "❌") || strings.HasPrefix(line, "⚡") {
			content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render("  "+line))
		} else if line == "" {
			content = append(content, line)
		} else {
			content = append(content, lipgloss.NewStyle().Foreground(colorComment).Render(line))
		}
	}

	content = append(content, "")

	// Show what will be applied
	monitor := m.pendingMonitor
	option := m.pendingScalingOption

	// Monitor info
	monitorTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("📱 Target Monitor")
	content = append(content, monitorTitle)
	content = append(content, "")

	monitorInfo := fmt.Sprintf("  %s (%dx%d@%.1fHz)",
		monitor.Name, monitor.Width, monitor.Height, monitor.RefreshRate)
	content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(monitorInfo))
	content = append(content, "")

	// Settings that will be applied
	settingsTitle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("🎯 Settings to Apply")
	content = append(content, settingsTitle)
	content = append(content, "")

	settings := []string{
		fmt.Sprintf("  Monitor Scale: %s", lipgloss.NewStyle().Foreground(colorGreen).Render(fmt.Sprintf("%.2fx", option.MonitorScale))),
		fmt.Sprintf("  GTK Scale: %s", lipgloss.NewStyle().Foreground(colorMagenta).Render(fmt.Sprintf("%dx", option.GTKScale))),
		fmt.Sprintf("  Font DPI: %s", lipgloss.NewStyle().Foreground(colorYellow).Render(fmt.Sprintf("%d", option.FontDPI))),
	}

	for _, setting := range settings {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(setting))
	}

	content = append(content, "")

	// Action name
	actionName := "Smart Scaling"
	if m.confirmationAction == ConfirmManualScaling {
		actionName = "Manual Scaling"
	}

	actionInfo := fmt.Sprintf("Action: %s - %s",
		lipgloss.NewStyle().Foreground(colorCyan).Render(actionName),
		lipgloss.NewStyle().Foreground(colorComment).Render(option.DisplayName))
	content = append(content, actionInfo)
	content = append(content, "")

	// Instructions
	instructionsStyle := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true)

	instructions := []string{
		"💡 Controls:",
		"  Enter/Space - Apply changes (refresh desktop)",
		"  Esc - Cancel and return",
	}

	for _, instruction := range instructions {
		content = append(content, instructionsStyle.Render(instruction))
	}

	// Simple, clean container - consistent with other screens
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorYellow).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderHelp(contentHeight int) string {
	var content []string

	// Award-winning help screen
	title := lipgloss.NewStyle().
		Foreground(colorYellow).
		Bold(true).
		Render("📖 Help & Controls")

	content = append(content, title)
	content = append(content, "")

	// Navigation section - clean and consistent
	navTitle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("🎮 Navigation")
	content = append(content, navTitle)
	content = append(content, "")

	navItems := []string{
		fmt.Sprintf("  %s %s   Navigate up/down in menus",
			lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("↑↓"),
			lipgloss.NewStyle().Foreground(colorGreen).Render("k j")),
		fmt.Sprintf("  %s %s   Navigate left/right (manual scaling)",
			lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("←→"),
			lipgloss.NewStyle().Foreground(colorCyan).Render("h l")),
		fmt.Sprintf("  %s       Select option or apply changes",
			lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("⏎")),
		fmt.Sprintf("  %s       Alternative selection",
			lipgloss.NewStyle().Foreground(colorBlue).Render("Space")),
	}

	for _, item := range navItems {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// Commands section - clean and consistent
	cmdTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("⌨️ Global Commands")
	content = append(content, cmdTitle)
	content = append(content, "")

	cmdItems := []string{
		fmt.Sprintf("  %s %s   Show this help screen",
			lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("h"),
			lipgloss.NewStyle().Foreground(colorYellow).Render("?")),
		fmt.Sprintf("  %s       Return to main menu",
			lipgloss.NewStyle().Foreground(colorMagenta).Bold(true).Render("Esc")),
		fmt.Sprintf("  %s       Quit application",
			lipgloss.NewStyle().Foreground(colorRed).Bold(true).Render("q")),
		fmt.Sprintf("  %s   Force quit",
			lipgloss.NewStyle().Foreground(colorRed).Bold(true).Render("Ctrl+C")),
	}

	for _, item := range cmdItems {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// Mode-specific controls
	modeTitle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("🎯 Mode-Specific Controls")
	content = append(content, modeTitle)
	content = append(content, "")

	modeItems := []string{
		fmt.Sprintf("  %s       Switch to manual scaling (from smart scaling)",
			lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("m")),
		fmt.Sprintf("  %s       Select control in manual scaling",
			lipgloss.NewStyle().Foreground(colorMagenta).Bold(true).Render("↑↓")),
		fmt.Sprintf("  %s       Adjust values in manual scaling",
			lipgloss.NewStyle().Foreground(colorMagenta).Bold(true).Render("←→")),
	}

	for _, item := range modeItems {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// About section - clean and minimal
	aboutTitle := lipgloss.NewStyle().Foreground(colorMagenta).Bold(true).Render("ℹ️ About")
	content = append(content, aboutTitle)
	content = append(content, "")

	aboutItems := []string{
		fmt.Sprintf("  Version: %s", lipgloss.NewStyle().Foreground(colorGreen).Render("1.0.0")),
		fmt.Sprintf("  Theme: %s", lipgloss.NewStyle().Foreground(colorMagenta).Render("Tokyo Night")),
		fmt.Sprintf("  Target: %s", lipgloss.NewStyle().Foreground(colorCyan).Render("Hyprland & Wayland")),
		fmt.Sprintf("  Built with: %s", lipgloss.NewStyle().Foreground(colorBlue).Render("Go + Bubbletea")),
	}

	for _, item := range aboutItems {
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

	// Footer message
	footer := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true).
		Render("💡 Press Esc to return to the main menu")
	content = append(content, footer)

	// Simple, clean container - consistent with other screens
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorYellow).
		Render(strings.Join(content, "\n"))
}
