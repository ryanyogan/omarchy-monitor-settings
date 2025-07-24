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
	ModeManualScaling
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

	// Scaling data
	scalingOptions     []ScalingOption
	selectedScalingOpt int

	// Manual scaling controls
	manualMonitorScale float64
	manualGTKScale     int
	manualFontDPI      int

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
			"Smart Scaling",
			"Manual Scaling",
			"Settings",
			"Help",
			"Exit",
		},
		isDemoMode: true, // Default to demo mode for testing

		// Initialize manual scaling defaults
		manualMonitorScale: 1.0,
		manualGTKScale:     1,
		manualFontDPI:      96,
	}

	// Initialize styles with Tokyo Night theme
	m.initStyles()

	// Load monitors using detection or demo data
	m.loadMonitors()

	// Load intelligent scaling options for the first monitor
	if len(m.monitors) > 0 {
		scalingManager := NewScalingManager()
		m.scalingOptions = scalingManager.GetIntelligentScalingOptions(m.monitors[0])
	}

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
				// Update scaling options when monitor changes
				if len(m.monitors) > 0 {
					scalingManager := NewScalingManager()
					m.scalingOptions = scalingManager.GetIntelligentScalingOptions(m.monitors[m.selectedMonitor])
					m.selectedScalingOpt = 0 // Reset to first option
				}
			}
		} else if m.mode == ModeScalingOptions {
			if m.selectedScalingOpt > 0 {
				m.selectedScalingOpt--
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
				// Update scaling options when monitor changes
				if len(m.monitors) > 0 {
					scalingManager := NewScalingManager()
					m.scalingOptions = scalingManager.GetIntelligentScalingOptions(m.monitors[m.selectedMonitor])
					m.selectedScalingOpt = 0 // Reset to first option
				}
			}
		} else if m.mode == ModeScalingOptions {
			if m.selectedScalingOpt < len(m.scalingOptions)-1 {
				m.selectedScalingOpt++
			}
		}

	case "enter", " ":
		if m.mode == ModeScalingOptions && len(m.scalingOptions) > 0 {
			// Apply the selected scaling option
			selectedOption := m.scalingOptions[m.selectedScalingOpt]
			monitor := m.monitors[m.selectedMonitor]

			configManager := NewConfigManager(m.isDemoMode)
			configManager.ApplyCompleteScalingOption(monitor, selectedOption)

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
		if m.mode == ModeManualScaling {
			m.mode = ModeScalingOptions
		} else {
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
	case ModeManualScaling:
		content = m.renderManualScaling(contentHeight)
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

	// Award-winning smart scaling screen
	title := lipgloss.NewStyle().
		Foreground(tokyoGreen).
		Bold(true).
		Render("🧠 Smart Scaling Recommendations")

	content = append(content, title)
	content = append(content, "")

	if len(m.monitors) > 0 {
		selectedMonitor := m.monitors[m.selectedMonitor]

		// Monitor info card with better formatting
		monitorHeader := fmt.Sprintf("%s • %s %s",
			lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render(selectedMonitor.Name),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(selectedMonitor.Make),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(selectedMonitor.Model))

		monitorSpecs := fmt.Sprintf("%dx%d @ %.0fHz",
			selectedMonitor.Width, selectedMonitor.Height, selectedMonitor.RefreshRate)

		monitorCard := lipgloss.NewStyle().
			Background(tokyoSurface).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(tokyoYellow).
			Render(fmt.Sprintf("%s\n%s", monitorHeader,
				lipgloss.NewStyle().Foreground(tokyoComment).Render(monitorSpecs)))

		content = append(content, monitorCard)
		content = append(content, "")

		// Smart recommendations section
		recTitle := lipgloss.NewStyle().Foreground(tokyoCyan).Bold(true).Render("🎯 Research-Based Options")
		content = append(content, recTitle)
		content = append(content, "")

		// Display scaling options
		optionColors := []lipgloss.Color{tokyoGreen, tokyoBlue, tokyoYellow, tokyoOrange}

		for i, option := range m.scalingOptions {
			if i >= len(optionColors) {
				break // Prevent overflow
			}

			color := optionColors[i%len(optionColors)]

			// Build option display
			var optionLines []string

			// Header with recommended badge
			headerText := option.DisplayName
			if option.IsRecommended {
				headerText += " " + lipgloss.NewStyle().
					Background(tokyoGreen).
					Foreground(tokyoBackground).
					Bold(true).
					Padding(0, 1).
					Render("RECOMMENDED")
			}

			if i == m.selectedScalingOpt {
				headerText = lipgloss.NewStyle().
					Foreground(color).
					Bold(true).
					Render("▶ " + headerText)
			} else {
				headerText = lipgloss.NewStyle().
					Foreground(color).
					Bold(true).
					Render("  " + headerText)
			}

			optionLines = append(optionLines, headerText)

			// Description with wrapping
			description := lipgloss.NewStyle().
				Foreground(tokyoSubtle).
				Width(m.width - 20). // Allow for wrapping
				Render("  " + option.Description)
			optionLines = append(optionLines, description)

			// Technical details
			details := fmt.Sprintf("  Monitor: %.1fx • GTK: %dx • Font DPI: %d • Result: %dx%d",
				option.MonitorScale, option.GTKScale, option.FontDPI,
				option.EffectiveWidth, option.EffectiveHeight)

			detailsStyled := lipgloss.NewStyle().
				Foreground(tokyoComment).
				Width(m.width - 20).
				Render(details)
			optionLines = append(optionLines, detailsStyled)

			// Reasoning with wrapping
			reasoning := lipgloss.NewStyle().
				Foreground(tokyoDark5).
				Italic(true).
				Width(m.width - 20).
				Render("  💡 " + option.Reasoning)
			optionLines = append(optionLines, reasoning)

			content = append(content, strings.Join(optionLines, "\n"))
			content = append(content, "") // Spacing between options
		}

		// Scaling explanations section
		content = append(content, "")
		explainTitle := lipgloss.NewStyle().Foreground(tokyoPurple).Bold(true).Render("📚 What Each Setting Does")
		content = append(content, explainTitle)

		configManager := NewConfigManager(m.isDemoMode)
		explanations := configManager.GetScalingExplanations()

		explainItems := []string{
			fmt.Sprintf("%s: %s",
				lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render("Monitor Scale"),
				lipgloss.NewStyle().Foreground(tokyoSubtle).Width(m.width-25).Render(explanations["monitor"])),
			fmt.Sprintf("%s: %s",
				lipgloss.NewStyle().Foreground(tokyoCyan).Bold(true).Render("GTK Scale"),
				lipgloss.NewStyle().Foreground(tokyoSubtle).Width(m.width-25).Render(explanations["gtk"])),
			fmt.Sprintf("%s: %s",
				lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render("Font DPI"),
				lipgloss.NewStyle().Foreground(tokyoSubtle).Width(m.width-25).Render(explanations["font"])),
		}

		explainSection := lipgloss.NewStyle().
			Background(tokyoSurface).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(tokyoPurple).
			Render(strings.Join(explainItems, "\n\n"))

		content = append(content, explainSection)

		if m.isDemoMode {
			content = append(content, "")
			demoNotice := lipgloss.NewStyle().
				Foreground(tokyoOrange).
				Background(tokyoSurface).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(tokyoOrange).
				Render("📱 Demo Mode: Changes will be simulated")
			content = append(content, demoNotice)
		}
	}

	// Instructions
	instructions := []string{
		lipgloss.NewStyle().Foreground(tokyoGreen).Render("↑↓") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" select option"),
		lipgloss.NewStyle().Foreground(tokyoBlue).Render("⏎") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" apply scaling"),
		lipgloss.NewStyle().Foreground(tokyoYellow).Render("m") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" manual mode"),
		lipgloss.NewStyle().Foreground(tokyoPurple).Render("esc") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" back"),
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
		BorderForeground(tokyoGreen).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderManualScaling(contentHeight int) string {
	var content []string

	// Award-winning manual scaling screen
	title := lipgloss.NewStyle().
		Foreground(tokyoPurple).
		Bold(true).
		Render("🔧 Manual Scaling Controls")

	content = append(content, title)
	content = append(content, "")

	if len(m.monitors) > 0 {
		selectedMonitor := m.monitors[m.selectedMonitor]

		// Monitor info card with current scaling
		monitorHeader := fmt.Sprintf("%s • %s %s",
			lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render(selectedMonitor.Name),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(selectedMonitor.Make),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(selectedMonitor.Model))

		monitorSpecs := fmt.Sprintf("%dx%d @ %.0fHz • Current Scale: %.2fx",
			selectedMonitor.Width, selectedMonitor.Height, selectedMonitor.RefreshRate, selectedMonitor.Scale)

		monitorCard := lipgloss.NewStyle().
			Background(tokyoSurface).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(tokyoYellow).
			Render(fmt.Sprintf("%s\n%s", monitorHeader,
				lipgloss.NewStyle().Foreground(tokyoComment).Render(monitorSpecs)))

		content = append(content, monitorCard)
		content = append(content, "")

		// Manual controls section
		controlsTitle := lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render("⚙️ Scaling Controls")
		content = append(content, controlsTitle)
		content = append(content, "")

		// Create three columns for each scaling type
		monitorCol := []string{
			lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render("Monitor Scale"),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render("(Compositor-level)"),
			"",
			lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render(fmt.Sprintf("%.2fx", m.manualMonitorScale)),
			"",
			lipgloss.NewStyle().Foreground(tokyoComment).Render("Range: 0.5x - 3.0x"),
			lipgloss.NewStyle().Foreground(tokyoComment).Render("Step: 0.25x"),
			"",
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("Scales everything"),
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("immediately. Works"),
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("with all apps."),
		}

		gtkCol := []string{
			lipgloss.NewStyle().Foreground(tokyoCyan).Bold(true).Render("GTK Scale"),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render("(Application-level)"),
			"",
			lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render(fmt.Sprintf("%dx", m.manualGTKScale)),
			"",
			lipgloss.NewStyle().Foreground(tokyoComment).Render("Range: 1x - 3x"),
			lipgloss.NewStyle().Foreground(tokyoComment).Render("Step: 1x (integer)"),
			"",
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("Scales GTK apps"),
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("(most Linux apps)."),
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("Requires logout."),
		}

		fontCol := []string{
			lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render("Font DPI"),
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render("(Text rendering)"),
			"",
			lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render(fmt.Sprintf("%d", m.manualFontDPI)),
			"",
			lipgloss.NewStyle().Foreground(tokyoComment).Render("Range: 72 - 288"),
			lipgloss.NewStyle().Foreground(tokyoComment).Render("Step: 12 DPI"),
			"",
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("Fine-grained text"),
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("scaling. Works with"),
			lipgloss.NewStyle().Foreground(tokyoDark5).Render("most applications."),
		}

		// Calculate effective resolution
		effectiveWidth := int(float64(selectedMonitor.Width) / m.manualMonitorScale)
		effectiveHeight := int(float64(selectedMonitor.Height) / m.manualMonitorScale)

		// Render columns side by side
		colWidth := (m.width - 12) / 3

		for i, col := range [][]string{monitorCol, gtkCol, fontCol} {
			for j, line := range col {
				if len(content) <= j+len(content)-len(monitorCol) {
					content = append(content, "")
				}
				// Style each column
				styledLine := lipgloss.NewStyle().
					Width(colWidth).
					Align(lipgloss.Left).
					Render(line)

				if i == 0 {
					content[len(content)-len(monitorCol)+j] = styledLine
				} else {
					content[len(content)-len(monitorCol)+j] += styledLine
				}
			}
		}

		// Results preview section
		content = append(content, "")
		content = append(content, "")
		resultsTitle := lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render("📊 Preview Results")
		content = append(content, resultsTitle)

		previewItems := []string{
			fmt.Sprintf("Effective Resolution: %s",
				lipgloss.NewStyle().Foreground(tokyoGreen).Bold(true).Render(fmt.Sprintf("%dx%d", effectiveWidth, effectiveHeight))),
			fmt.Sprintf("Screen Real Estate: %s",
				lipgloss.NewStyle().Foreground(tokyoBlue).Bold(true).Render(fmt.Sprintf("%.0f%%", 100.0/m.manualMonitorScale))),
			fmt.Sprintf("Font DPI Multiplier: %s",
				lipgloss.NewStyle().Foreground(tokyoYellow).Bold(true).Render(fmt.Sprintf("%.1fx", float64(m.manualFontDPI)/96.0))),
		}

		previewSection := lipgloss.NewStyle().
			Background(tokyoSurface).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(tokyoGreen).
			Render(strings.Join(previewItems, "\n"))

		content = append(content, previewSection)

		if m.isDemoMode {
			content = append(content, "")
			demoNotice := lipgloss.NewStyle().
				Foreground(tokyoOrange).
				Background(tokyoSurface).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(tokyoOrange).
				Render("📱 Demo Mode: Use ⏎ to preview changes")
			content = append(content, demoNotice)
		}
	}

	// Enhanced instructions
	instructions := []string{
		lipgloss.NewStyle().Foreground(tokyoBlue).Render("1-3") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" select control"),
		lipgloss.NewStyle().Foreground(tokyoGreen).Render("↑↓") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" adjust value"),
		lipgloss.NewStyle().Foreground(tokyoYellow).Render("⏎") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" apply changes"),
		lipgloss.NewStyle().Foreground(tokyoPurple).Render("esc") +
			lipgloss.NewStyle().Foreground(tokyoSubtle).Render(" back to smart scaling"),
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
		BorderForeground(tokyoPurple).
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
