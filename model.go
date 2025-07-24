package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tokyo Night Colors - Full vibrant palette
var (
	// Background colors
	tokyoBackground = lipgloss.Color("#1a1b26")
	tokyoSurface    = lipgloss.Color("#24283b")
	tokyoSidebar    = lipgloss.Color("#16161e")
	tokyoFloat      = lipgloss.Color("#1d202f")

	// Foreground colors
	tokyoForeground = lipgloss.Color("#c0caf5")
	tokyoComment    = lipgloss.Color("#565f89")
	tokyoSubtle     = lipgloss.Color("#9aa5ce")
	tokyoDark3      = lipgloss.Color("#545c7e")
	tokyoDark5      = lipgloss.Color("#737aa2")

	// Accent colors - Full vibrant palette
	tokyoBlue    = lipgloss.Color("#7aa2f7")
	tokyoCyan    = lipgloss.Color("#7dcfff")
	tokyoGreen   = lipgloss.Color("#9ece6a")
	tokyoYellow  = lipgloss.Color("#e0af68")
	tokyoOrange  = lipgloss.Color("#ff9e64")
	tokyoRed     = lipgloss.Color("#f7768e")
	tokyoPurple  = lipgloss.Color("#bb9af7")
	tokyoMagenta = lipgloss.Color("#c0a8f7")
	tokyoTeal    = lipgloss.Color("#1abc9c")
	tokyoPink    = lipgloss.Color("#f7768e")

	// Special Tokyo Night effects
	tokyoBlue0 = lipgloss.Color("#3d59a1")
	tokyoBlue1 = lipgloss.Color("#2ac3de")
	tokyoBlue2 = lipgloss.Color("#0db9d7")
	tokyoBlue6 = lipgloss.Color("#b4f9f8")
	tokyoBlue7 = lipgloss.Color("#394b70")
)

// AppMode represents different screens/modes in the TUI
type AppMode int

const (
	ModeDashboard AppMode = iota
	ModeMonitorSelection
	ModeScalingOptions
	ModeSettings
	ModeHelp
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

	// Monitor data
	monitors        []Monitor
	selectedMonitor int
	isDemoMode      bool

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

// NewModel creates and initializes a new Model
func NewModel() Model {
	m := Model{
		mode:           ModeDashboard,
		selectedOption: 0,
		menuItems: []string{
			"Dashboard",
			"Monitor Selection",
			"Scaling Options",
			"Settings",
			"Help",
			"Exit",
		},
		isDemoMode: true, // Default to demo mode for testing
	}

	// Initialize styles with Tokyo Night theme
	m.initStyles()

	// Load monitors using detection or demo data
	m.loadMonitors()

	return m
}

func (m *Model) initStyles() {
	// Header style - cleaner btop-like header
	m.headerStyle = lipgloss.NewStyle().
		Background(tokyoSurface).
		Foreground(tokyoForeground).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoBlue)

	// Footer style
	m.footerStyle = lipgloss.NewStyle().
		Background(tokyoSidebar).
		Foreground(tokyoComment).
		Padding(0, 1)

	// Title style
	m.titleStyle = lipgloss.NewStyle().
		Foreground(tokyoBlue).
		Bold(true).
		Underline(true)

		// Selected menu item style - no background, just text styling
	m.selectedStyle = lipgloss.NewStyle().
		Foreground(tokyoBlue).
		Bold(true)

	// Unselected menu item style
	m.unselectedStyle = lipgloss.NewStyle().
		Foreground(tokyoForeground)

	// Help text style
	m.helpStyle = lipgloss.NewStyle().
		Foreground(tokyoComment).
		Italic(true)

	// Error style
	m.errorStyle = lipgloss.NewStyle().
		Foreground(tokyoRed).
		Bold(true)

	// Success style
	m.successStyle = lipgloss.NewStyle().
		Foreground(tokyoGreen).
		Bold(true)
}

func (m *Model) loadMonitors() {
	detector := NewMonitorDetector()
	monitors, err := detector.DetectMonitors()

	if debugMode {
		fmt.Printf("DEBUG: DetectMonitors returned %d monitors, error: %v\n", len(monitors), err)
	}

	if err != nil {
		// Fallback to demo monitors
		if debugMode {
			fmt.Printf("DEBUG: Setting demo mode due to detection error\n")
		}
		m.isDemoMode = true
		monitors, _ = detector.getFallbackMonitors()
	} else {
		if debugMode {
			fmt.Printf("DEBUG: Setting live mode - detected real monitors\n")
		}
		m.isDemoMode = false
	}

	// Override demo mode if force-live flag is set
	if forceLiveMode {
		if debugMode {
			fmt.Printf("DEBUG: Force-live mode enabled, overriding demo mode\n")
		}
		m.isDemoMode = false
	}

	if debugMode {
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
		if m.mode == ModeDashboard {
			if m.selectedOption > 0 {
				m.selectedOption--
			}
		} else if m.mode == ModeMonitorSelection {
			if m.selectedMonitor > 0 {
				m.selectedMonitor--
			}
		}

	case "down", "j":
		if m.mode == ModeDashboard {
			if m.selectedOption < len(m.menuItems)-1 {
				m.selectedOption++
			}
		} else if m.mode == ModeMonitorSelection {
			if m.selectedMonitor < len(m.monitors)-1 {
				m.selectedMonitor++
			}
		}

	case "enter", " ":
		return m.handleSelection()

	case "h", "?":
		m.mode = ModeHelp

	case "esc":
		m.mode = ModeDashboard
		// Reset selection when returning to dashboard
		m.selectedOption = 0
	}

	return m, nil
}

func (m Model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.selectedOption {
	case 0: // Dashboard
		m.mode = ModeDashboard
	case 1: // Monitor Selection
		m.mode = ModeMonitorSelection
	case 2: // Scaling Options
		m.mode = ModeScalingOptions
	case 3: // Settings
		m.mode = ModeSettings
	case 4: // Help
		m.mode = ModeHelp
	case 5: // Exit
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
			Foreground(tokyoRed).
			Render("Terminal too small\nPlease resize to at least 80x20")
	}

	if !m.ready {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(tokyoBlue).
			Render("Initializing stunning TUI...")
	}

	// Calculate precise dimensions for award-winning layout
	headerHeight := 4                                           // Slightly more breathing room
	footerHeight := 2                                           // Better footer spacing
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
	case ModeSettings:
		content = m.renderSettings(contentHeight)
	case ModeHelp:
		content = m.renderHelp(contentHeight)
	default:
		content = m.renderDashboard(contentHeight)
	}

	// Create award-winning layout with perfect spacing
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
	// Award-winning header design with perfect typography
	title := lipgloss.NewStyle().
		Foreground(tokyoBlue).
		Bold(true).
		Render("Hyprland Monitor Manager")

	var statusBadge string
	if m.isDemoMode {
		statusBadge = lipgloss.NewStyle().
			Background(tokyoOrange).
			Foreground(tokyoBackground).
			Bold(true).
			Padding(0, 1).
			Render(" DEMO ")
	} else {
		statusBadge = lipgloss.NewStyle().
			Background(tokyoGreen).
			Foreground(tokyoBackground).
			Bold(true).
			Padding(0, 1).
			Render(" LIVE ")
	}

	// Professional header layout
	headerLeft := lipgloss.JoinVertical(lipgloss.Left,
		title,
		lipgloss.NewStyle().Foreground(tokyoSubtle).Render("Beautiful Display Configuration"),
	)

	headerRight := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Render(statusBadge)

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerLeft,
		lipgloss.NewStyle().Width(m.width-40).Render(""), // Spacer
		headerRight,
	)

	return lipgloss.NewStyle().
		Width(m.width-2).
		Background(tokyoSurface).
		Foreground(tokyoForeground).
		Padding(1, 2).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoBlue).
		Render(headerContent)
}

func (m Model) renderFooter() string {
	// Award-winning footer with elegant key hints
	keyStyle := lipgloss.NewStyle().
		Background(tokyoFloat).
		Foreground(tokyoBlue).
		Bold(true).
		Padding(0, 1).
		Margin(0, 1)

	textStyle := lipgloss.NewStyle().
		Foreground(tokyoSubtle)

	controls := []string{
		keyStyle.Copy().Foreground(tokyoGreen).Render("↑↓") + textStyle.Render("navigate"),
		keyStyle.Copy().Foreground(tokyoBlue).Render("⏎") + textStyle.Render("select"),
		keyStyle.Copy().Foreground(tokyoYellow).Render("h") + textStyle.Render("help"),
		keyStyle.Copy().Foreground(tokyoPurple).Render("esc") + textStyle.Render("back"),
		keyStyle.Copy().Foreground(tokyoRed).Render("q") + textStyle.Render("quit"),
	}

	helpText := strings.Join(controls, "  ")

	return lipgloss.NewStyle().
		Width(m.width-2).
		Background(tokyoSidebar).
		Align(lipgloss.Center).
		Padding(1, 2).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoComment).
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
			Foreground(tokyoBlue).
			Bold(true).
			Render("Navigation"),
	)
	leftPanel = append(leftPanel, "")

	menuColors := []lipgloss.Color{tokyoBlue, tokyoCyan, tokyoGreen, tokyoYellow, tokyoOrange, tokyoRed}

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
				Foreground(tokyoSubtle).
				Render(item)
			line = fmt.Sprintf("  %s", text)
		}
		leftPanel = append(leftPanel, line)
		leftPanel = append(leftPanel, "") // Breathing room
	}

	// Elegant monitor panel
	rightPanel = append(rightPanel,
		lipgloss.NewStyle().
			Foreground(tokyoCyan).
			Bold(true).
			Render("Display Overview"),
	)
	rightPanel = append(rightPanel, "")

	monitorColors := []lipgloss.Color{tokyoGreen, tokyoBlue, tokyoYellow, tokyoMagenta}

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
				statusStyle = lipgloss.NewStyle().Foreground(tokyoGreen)
			} else {
				statusIcon = "○"
				statusStyle = lipgloss.NewStyle().Foreground(tokyoBlue)
			}
		} else {
			statusIcon = "◦"
			statusStyle = lipgloss.NewStyle().Foreground(tokyoComment)
		}

		// Beautiful monitor card
		header := fmt.Sprintf("%s %s",
			statusStyle.Render(statusIcon),
			lipgloss.NewStyle().Foreground(color).Bold(true).Render(monitor.Name),
		)

		details := []string{
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(fmt.Sprintf("  %s %s", monitor.Make, monitor.Model)),
			lipgloss.NewStyle().Foreground(tokyoComment).Render(fmt.Sprintf("  %dx%d @ %.0fHz", monitor.Width, monitor.Height, monitor.RefreshRate)),
			lipgloss.NewStyle().Foreground(tokyoComment).Render(fmt.Sprintf("  Scale: %.1fx", monitor.Scale)),
		}

		rightPanel = append(rightPanel, header)
		rightPanel = append(rightPanel, details...)
		rightPanel = append(rightPanel, "") // Card spacing
	}

	if m.isDemoMode {
		demoNotice := lipgloss.NewStyle().
			Foreground(tokyoOrange).
			Italic(true).
			Render("📱 Demo Mode Active")
		rightPanel = append(rightPanel, demoNotice)
	}

	// Award-winning panel styling
	leftContent := lipgloss.NewStyle().
		Width(leftWidth).
		Height(contentHeight - 2).
		Padding(2).
		Background(tokyoFloat).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoBlue).
		Render(strings.Join(leftPanel, "\n"))

	rightContent := lipgloss.NewStyle().
		Width(rightWidth).
		Height(contentHeight - 2).
		Padding(2).
		Background(tokyoFloat).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoCyan).
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
		Foreground(tokyoBlue).
		Bold(true).
		Render("Monitor Selection")

	subtitle := lipgloss.NewStyle().
		Foreground(tokyoSubtle).
		Render("Choose a display to configure")

	content = append(content, title)
	content = append(content, subtitle)
	content = append(content, "")

	monitorColors := []lipgloss.Color{tokyoGreen, tokyoBlue, tokyoYellow, tokyoMagenta}

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
				statusStyle = lipgloss.NewStyle().Background(tokyoGreen).Foreground(tokyoBackground).Bold(true).Padding(0, 1)
			} else {
				statusText = "ACTIVE"
				statusStyle = lipgloss.NewStyle().Background(tokyoBlue).Foreground(tokyoBackground).Bold(true).Padding(0, 1)
			}
		} else {
			statusText = "INACTIVE"
			statusStyle = lipgloss.NewStyle().Background(tokyoComment).Foreground(tokyoBackground).Bold(true).Padding(0, 1)
		}

		// Beautiful monitor card
		nameStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
		detailStyle := lipgloss.NewStyle().Foreground(tokyoSubtle)

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
		lipgloss.NewStyle().Foreground(tokyoYellow).Render("⏎") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" Configure selected monitor"),
		lipgloss.NewStyle().Foreground(tokyoPurple).Render("esc") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" Return to main menu"),
	}

	content = append(content, "")
	content = append(content, strings.Join(instructions, "  "))

	// Beautiful container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(tokyoFloat).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoBlue).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderScalingOptions(contentHeight int) string {
	var content []string

	// Award-winning scaling options screen
	title := lipgloss.NewStyle().
		Foreground(tokyoGreen).
		Bold(true).
		Render("Scaling Configuration")

	content = append(content, title)
	content = append(content, "")

	if len(m.monitors) > 0 {
		selectedMonitor := m.monitors[m.selectedMonitor]

		// Monitor info card
		monitorCard := lipgloss.NewStyle().
			Background(tokyoSurface).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(tokyoYellow).
			Render(fmt.Sprintf("%s %s\n%dx%d @ %.0fHz",
				lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render(selectedMonitor.Name),
				lipgloss.NewStyle().Foreground(tokyoSubtle).Render(selectedMonitor.Make+" "+selectedMonitor.Model),
				selectedMonitor.Width, selectedMonitor.Height, selectedMonitor.RefreshRate))

		content = append(content, monitorCard)
		content = append(content, "")

		// Get scaling recommendation
		scalingManager := NewScalingManager()
		recommendation := scalingManager.GetRecommendedScale(selectedMonitor)

		// Recommendations section
		recTitle := lipgloss.NewStyle().Foreground(tokyoCyan).Bold(true).Render("✨ Smart Recommendations")
		content = append(content, recTitle)
		content = append(content, "")

		recItems := []string{
			fmt.Sprintf("Monitor Scale: %s",
				lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render(fmt.Sprintf("%.1fx", recommendation.MonitorScale))),
			fmt.Sprintf("Font Scale: %s",
				lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render(fmt.Sprintf("%.1fx", recommendation.FontScale))),
			fmt.Sprintf("Effective Resolution: %s",
				lipgloss.NewStyle().Foreground(tokyoPurple).Bold(true).Render(fmt.Sprintf("%dx%d", recommendation.EffectiveWidth, recommendation.EffectiveHeight))),
		}

		for _, item := range recItems {
			content = append(content, "  "+item)
		}

		content = append(content, "")
		reasoning := lipgloss.NewStyle().
			Foreground(tokyoSubtle).
			Italic(true).
			Render("💡 " + recommendation.Reasoning)
		content = append(content, reasoning)
		content = append(content, "")

		// Available options
		optionsTitle := lipgloss.NewStyle().Foreground(tokyoOrange).Bold(true).Render("Available Options")
		content = append(content, optionsTitle)
		content = append(content, "")

		scaleOptions := []struct {
			scale string
			desc  string
			color lipgloss.Color
		}{
			{"1.0x", "Native resolution", tokyoBlue},
			{"1.5x", "150% scaling", tokyoCyan},
			{"2.0x", "200% scaling (4K recommended)", tokyoGreen},
		}

		for _, option := range scaleOptions {
			scaleText := lipgloss.NewStyle().Foreground(option.color).Bold(true).Render(option.scale)
			line := fmt.Sprintf("  %s - %s", scaleText, option.desc)
			content = append(content, line)
		}

		if m.isDemoMode {
			content = append(content, "")
			demoNotice := lipgloss.NewStyle().
				Foreground(tokyoOrange).
				Background(tokyoSurface).
				Padding(0, 1).
				Render("📱 Demo Mode: Changes won't be applied")
			content = append(content, demoNotice)
		}
	}

	// Beautiful container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(tokyoFloat).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoGreen).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderSettings(contentHeight int) string {
	var content []string

	// Award-winning settings screen
	title := lipgloss.NewStyle().
		Foreground(tokyoPurple).
		Bold(true).
		Render("Application Settings")

	content = append(content, title)
	content = append(content, "")

	// Application info section
	appSection := lipgloss.NewStyle().
		Background(tokyoSurface).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoBlue).
		Render(fmt.Sprintf("%s\n%s\n%s\n%s",
			lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render("📱 Application Info"),
			fmt.Sprintf("Version: %s", lipgloss.NewStyle().Foreground(tokyoGreen).Render("0.1.0")),
			fmt.Sprintf("Theme: %s", lipgloss.NewStyle().Foreground(tokyoPurple).Render("Tokyo Night")),
			fmt.Sprintf("Mode: %s", func() string {
				if m.isDemoMode {
					return lipgloss.NewStyle().Foreground(tokyoOrange).Render("Demo")
				}
				return lipgloss.NewStyle().Foreground(tokyoGreen).Render("Live")
			}()),
		))

	content = append(content, appSection)
	content = append(content, "")

	// Detection methods section
	detector := NewMonitorDetector()
	var detectionItems []string

	detectionItems = append(detectionItems, lipgloss.NewStyle().Foreground(tokyoCyan).Bold(true).Render("🔍 Detection Methods"))

	methods := []struct {
		name string
		cmd  string
	}{
		{"Hyprctl", "hyprctl"},
		{"wlr-randr", "wlr-randr"},
		{"xrandr", "xrandr"},
	}

	for _, method := range methods {
		var status string
		if detector.commandExists(method.cmd) {
			status = lipgloss.NewStyle().Foreground(tokyoGreen).Render("✓ Available")
		} else {
			status = lipgloss.NewStyle().Foreground(tokyoRed).Render("✗ Not found")
		}
		detectionItems = append(detectionItems, fmt.Sprintf("%s: %s", method.name, status))
	}

	detectionSection := lipgloss.NewStyle().
		Background(tokyoSurface).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoCyan).
		Render(strings.Join(detectionItems, "\n"))

	content = append(content, detectionSection)
	content = append(content, "")

	// Configuration section
	configItems := []string{
		lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render("⚙️ Configuration"),
		fmt.Sprintf("Target: %s", lipgloss.NewStyle().Foreground(tokyoGreen).Render("Hyprland + Wayland")),
		fmt.Sprintf("Fallbacks: %s", lipgloss.NewStyle().Foreground(tokyoBlue).Render("wlr-randr, xrandr")),
		fmt.Sprintf("Font Scaling: %s", lipgloss.NewStyle().Foreground(tokyoPurple).Render("GTK, Alacritty, Neovim")),
	}

	configSection := lipgloss.NewStyle().
		Background(tokyoSurface).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoYellow).
		Render(strings.Join(configItems, "\n"))

	content = append(content, configSection)

	// Beautiful container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(tokyoFloat).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoPurple).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderHelp(contentHeight int) string {

	// Award-winning help screen
	title := lipgloss.NewStyle().
		Foreground(tokyoYellow).
		Bold(true).
		Render("Help & Controls")

	// Navigation section
	navItems := []string{
		lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render("🎮 Navigation"),
		fmt.Sprintf("%s / %s   Move up/down",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoGreen).Padding(0, 1).Render("↑"),
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoGreen).Padding(0, 1).Render("k")),
		fmt.Sprintf("%s / %s   Move down/up",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoGreen).Padding(0, 1).Render("↓"),
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoGreen).Padding(0, 1).Render("j")),
		fmt.Sprintf("%s      Select option",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoBlue).Padding(0, 1).Render("Enter")),
		fmt.Sprintf("%s      Alternative select",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoBlue).Padding(0, 1).Render("Space")),
	}

	navSection := lipgloss.NewStyle().
		Background(tokyoSurface).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoGreen).
		Render(strings.Join(navItems, "\n"))

	// Commands section
	cmdItems := []string{
		lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render("⌨️  Global Commands"),
		fmt.Sprintf("%s / %s   Show this help",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoYellow).Padding(0, 1).Render("h"),
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoYellow).Padding(0, 1).Render("?")),
		fmt.Sprintf("%s      Return to main menu",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoPurple).Padding(0, 1).Render("Esc")),
		fmt.Sprintf("%s        Quit application",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoRed).Padding(0, 1).Render("q")),
		fmt.Sprintf("%s   Force quit",
			lipgloss.NewStyle().Background(tokyoFloat).Foreground(tokyoRed).Padding(0, 1).Render("Ctrl+C")),
	}

	cmdSection := lipgloss.NewStyle().
		Background(tokyoSurface).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoBlue).
		Render(strings.Join(cmdItems, "\n"))

	// About section
	aboutItems := []string{
		lipgloss.NewStyle().Foreground(tokyoPurple).Bold(true).Render("ℹ️  About"),
		fmt.Sprintf("Theme: %s", lipgloss.NewStyle().Foreground(tokyoPurple).Render("Tokyo Night")),
		fmt.Sprintf("Target: %s", lipgloss.NewStyle().Foreground(tokyoCyan).Render("Hyprland & Wayland")),
		fmt.Sprintf("Version: %s", lipgloss.NewStyle().Foreground(tokyoGreen).Render("0.1.0")),
		fmt.Sprintf("Built with: %s", lipgloss.NewStyle().Foreground(tokyoBlue).Render("Go + Bubbletea")),
	}

	aboutSection := lipgloss.NewStyle().
		Background(tokyoSurface).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoPurple).
		Render(strings.Join(aboutItems, "\n"))

	// Combine sections
	content := []string{
		title,
		"",
		navSection,
		"",
		cmdSection,
		"",
		aboutSection,
		"",
		lipgloss.NewStyle().
			Foreground(tokyoSubtle).
			Italic(true).
			Render("Press Esc to return to the main menu"),
	}

	// Beautiful container
	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(tokyoFloat).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tokyoYellow).
		Render(strings.Join(content, "\n"))
}
