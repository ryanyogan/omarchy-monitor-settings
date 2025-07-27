package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/ryanyogan/omarchy-monitor-settings/internal/app"
	"github.com/ryanyogan/omarchy-monitor-settings/internal/monitor"
	"github.com/ryanyogan/omarchy-monitor-settings/pkg/types"
	"github.com/ryanyogan/omarchy-monitor-settings/pkg/utils"
)

var (
	colorBackground = lipgloss.Color("")
	colorSurface    = lipgloss.Color("0")
	colorFloat      = lipgloss.Color("8")

	colorForeground = lipgloss.Color("")
	colorComment    = lipgloss.Color("8")
	colorSubtle     = lipgloss.Color("7")

	colorBlue    = lipgloss.Color("4")
	colorCyan    = lipgloss.Color("6")
	colorGreen   = lipgloss.Color("2")
	colorYellow  = lipgloss.Color("3")
	colorRed     = lipgloss.Color("1")
	colorMagenta = lipgloss.Color("5")
)

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

	var theme string
	if isDark {
		theme = "Dark"
	} else {
		theme = "Light"
	}

	return fmt.Sprintf("Terminal Adaptive (%s, %s)", profileName, theme)
}

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

type ConfirmationAction int

const (
	ConfirmNone ConfirmationAction = iota
	ConfirmSmartScaling
	ConfirmManualScaling
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

type Model struct {
	width  int
	height int

	mode AppMode

	menuItems      []string
	selectedOption int

	monitors           []monitor.Monitor
	selectedMonitor    int
	selectedScalingOpt int
	scalingOptions     []monitor.ScalingOption

	manualMonitorScale    float64
	manualGTKScale        int
	manualFontDPI         int
	selectedManualControl int

	confirmationAction ConfirmationAction
	pendingOption      monitor.ScalingOption
	pendingMonitor     monitor.Monitor

	isDemoMode bool
	ready      bool

	services *app.Services

	cachedTerminalTheme string
	cachedCommandStatus map[string]bool

	headerStyle     lipgloss.Style
	footerStyle     lipgloss.Style
	titleStyle      lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	helpStyle       lipgloss.Style
	errorStyle      lipgloss.Style
	successStyle    lipgloss.Style
}

func NewModel() Model {
	config := &app.Config{
		NoHyprlandCheck: true,
		DebugMode:       false,
		ForceLiveMode:   false,
		IsTestMode:      true,
	}
	services := app.NewServices(config)
	return NewModelWithServices(services)
}

func NewModelWithServices(services *app.Services) Model {
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
		isDemoMode: true,
		services:   services,

		manualMonitorScale:    1.0,
		manualGTKScale:        1,
		manualFontDPI:         types.BaseDPI,
		selectedManualControl: 0,

		cachedCommandStatus: make(map[string]bool),
	}

	m.initStyles()

	m.cachedTerminalTheme = getTerminalThemeInfo()

	commands := []string{"hyprctl", "wlr-randr"}
	for _, cmd := range commands {
		m.cachedCommandStatus[cmd] = utils.CommandExists(cmd)
	}

	m.loadMonitors()

	if len(m.monitors) > 0 {
		m.scalingOptions = services.ScalingManager.GetIntelligentScalingOptions(m.monitors[0])
	}

	return m
}

func (m *Model) initStyles() {
	m.headerStyle = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorForeground).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue)

	m.footerStyle = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorComment).
		Padding(0, 1)

	m.titleStyle = lipgloss.NewStyle().
		Foreground(colorBlue).
		Bold(true).
		Underline(true)

	m.selectedStyle = lipgloss.NewStyle().
		Foreground(colorBlue).
		Bold(true)

	m.unselectedStyle = lipgloss.NewStyle().
		Foreground(colorForeground)

	m.helpStyle = lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true)

	m.errorStyle = lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)

	m.successStyle = lipgloss.NewStyle().
		Foreground(colorGreen).
		Bold(true)
}

func (m *Model) loadMonitors() {
	monitors, err := m.services.MonitorDetector.DetectMonitors()

	if m.services.Config.DebugMode {
		fmt.Printf("DEBUG: DetectMonitors returned %d monitors, error: %v\n", len(monitors), err)
	}

	if err != nil {
		if m.services.Config.DebugMode {
			fmt.Printf("DEBUG: Setting demo mode due to detection error\n")
		}
		m.isDemoMode = true

		if detector, ok := m.services.MonitorDetector.(*monitor.Detector); ok {
			monitors = detector.GetFallbackMonitors()
		}
	} else {
		if m.services.Config.DebugMode {
			fmt.Printf("DEBUG: Setting live mode - detected real monitors\n")
		}
		m.isDemoMode = false
	}

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

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

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
				if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
					m.scalingOptions = m.services.ScalingManager.GetIntelligentScalingOptions(m.monitors[m.selectedMonitor])
					m.selectedScalingOpt = 0
				}
			}
		case ModeScalingOptions:
			if m.selectedScalingOpt > 0 {
				m.selectedScalingOpt--
			}
		case ModeManualScaling:
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
				if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
					m.scalingOptions = m.services.ScalingManager.GetIntelligentScalingOptions(m.monitors[m.selectedMonitor])
					m.selectedScalingOpt = 0
				}
			}
		case ModeScalingOptions:
			if m.selectedScalingOpt < len(m.scalingOptions)-1 {
				m.selectedScalingOpt++
			}
		case ModeManualScaling:
			if m.selectedManualControl < 2 {
				m.selectedManualControl++
			}
		}

	case "left":
		if m.mode == ModeManualScaling {
			switch m.selectedManualControl {
			case 0:
				m.manualMonitorScale = utils.FindNextValidScale(m.manualMonitorScale, false, types.ValidHyprlandScales)
			case 1:
				if m.manualGTKScale > types.MinGTKScale {
					m.manualGTKScale--
				}
			case 2:
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
			switch m.selectedManualControl {
			case 0:
				m.manualMonitorScale = utils.FindNextValidScale(m.manualMonitorScale, true, types.ValidHyprlandScales)
			case 1:
				if m.manualGTKScale < types.MaxGTKScale {
					m.manualGTKScale++
				}
			case 2:
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
			m.mode = ModeDashboard
			m.selectedOption = 0
			return m, nil
		} else if m.mode == ModeScalingOptions && len(m.scalingOptions) > 0 && m.selectedScalingOpt < len(m.scalingOptions) {
			selectedOption := m.scalingOptions[m.selectedScalingOpt]
			if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
				m.confirmationAction = ConfirmSmartScaling
				m.pendingOption = selectedOption
				m.pendingMonitor = m.monitors[m.selectedMonitor]
				m.mode = ModeConfirmation
			}
			return m, nil
		} else if m.mode == ModeManualScaling {
			if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
				m.confirmationAction = ConfirmManualScaling
				m.pendingMonitor = m.monitors[m.selectedMonitor]
				m.pendingOption = monitor.ScalingOption{
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
			switch m.confirmationAction {
			case ConfirmSmartScaling:
				_ = m.services.ConfigManager.ApplyCompleteScalingOption(m.pendingMonitor, m.pendingOption)
				// Update the local monitor data with new scale
				if m.selectedMonitor < len(m.monitors) {
					m.monitors[m.selectedMonitor].Scale = m.pendingOption.MonitorScale
				}
			case ConfirmManualScaling:
				_ = m.services.ConfigManager.ApplyMonitorScale(m.pendingMonitor, m.manualMonitorScale)
				_ = m.services.ConfigManager.ApplyGTKScale(m.manualGTKScale)
				_ = m.services.ConfigManager.ApplyFontDPI(m.manualFontDPI)
				// Update the local monitor data with new scale
				if m.selectedMonitor < len(m.monitors) {
					m.monitors[m.selectedMonitor].Scale = m.manualMonitorScale
				}
			}
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
			m.selectedOption = 0
		case ModeConfirmation:
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
			m.selectedOption = 0
		}
	}

	return m, nil
}

func (m Model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.selectedOption {
	case 0:
		m.mode = ModeDashboard
	case 1:
		m.mode = ModeMonitorSelection
	case 2:
		m.mode = ModeScalingOptions
	case 3:
		m.mode = ModeManualScaling
	case 4:
		m.mode = ModeSettings
	case 5:
		m.mode = ModeHelp
	case 6:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	if m.width < types.MinTerminalWidth || m.height < types.MinTerminalHeight {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorRed).
			Render(types.ErrTerminalTooSmall)
	}

	if !m.ready {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorBlue).
			Render("Initializing stunning TUI...")
	}

	headerHeight := 7
	footerHeight := 2
	contentHeight := m.height - headerHeight - footerHeight - 2

	if contentHeight < 10 {
		contentHeight = 10
	}

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

	header := m.renderHeader()
	footer := m.renderFooter()

	styledContent := lipgloss.NewStyle().
		Width(m.width-4).
		Height(contentHeight).
		Margin(1, 2).
		Render(content)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		styledContent,
		footer,
	)
}

func (m Model) renderHeader() string {
	availableWidth := m.width - 8
	leftWidth := availableWidth * 2 / 5
	rightWidth := availableWidth - leftWidth - 4

	if leftWidth < 25 {
		leftWidth = 25
	}
	if rightWidth < 30 {
		rightWidth = 30
	}

	totalHeaderWidth := leftWidth + 2 + rightWidth + 2

	return lipgloss.NewStyle().
		Width(totalHeaderWidth).
		Background(colorBackground).
		Foreground(colorBlue).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 2).
		Margin(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorComment).
		Render("Display Settings")
}

func (m Model) renderFooter() string {
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

	availableWidth := m.width - 8
	leftWidth := availableWidth * 2 / 5
	rightWidth := availableWidth - leftWidth - 4

	if leftWidth < 25 {
		leftWidth = 25
	}
	if rightWidth < 30 {
		rightWidth = 30
	}

	totalFooterWidth := leftWidth + 2 + rightWidth + 2

	return lipgloss.NewStyle().
		Width(totalFooterWidth).
		Background(colorBackground).
		Align(lipgloss.Center).
		Padding(1, 2).
		Margin(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorComment).
		Render(helpText)
}

func (m Model) renderDashboard(contentHeight int) string {
	availableWidth := m.width - 8
	leftWidth := availableWidth * 2 / 5
	rightWidth := availableWidth - leftWidth - 4

	if leftWidth < 25 {
		leftWidth = 25
	}
	if rightWidth < 30 {
		rightWidth = 30
	}

	var leftPanel []string
	var rightPanel []string

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
		leftPanel = append(leftPanel, "")
	}

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
			break
		}

		color := monitorColors[i]

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

		header := fmt.Sprintf("%s %s",
			statusStyle.Render(statusIcon),
			lipgloss.NewStyle().Foreground(color).Bold(true).Render(monitor.Name),
		)

		if i == m.selectedMonitor {
			selectedIndicator := lipgloss.NewStyle().
				Foreground(colorYellow).
				Bold(true).
				Render(" 👆 CURRENT")
			header = header + selectedIndicator
		}

		details := []string{
			lipgloss.NewStyle().Foreground(colorSubtle).Render(fmt.Sprintf("  %s %s", monitor.Make, monitor.Model)),
			lipgloss.NewStyle().Foreground(colorComment).Render(fmt.Sprintf("  %s @ %.0fHz", utils.FormatResolution(monitor.Width, monitor.Height), monitor.RefreshRate)),
			lipgloss.NewStyle().Foreground(colorComment).Render(fmt.Sprintf("  Scale: %.1fx", monitor.Scale)),
		}

		if i == m.selectedMonitor {
			details = append(details, lipgloss.NewStyle().
				Foreground(colorYellow).
				Italic(true).
				Render("  → Scaling changes will apply here"))
		}

		rightPanel = append(rightPanel, header)
		rightPanel = append(rightPanel, details...)
		rightPanel = append(rightPanel, "")
	}

	if m.isDemoMode {
		demoNotice := lipgloss.NewStyle().
			Foreground(colorYellow).
			Italic(true).
			Render("📱 Demo Mode Active")
		rightPanel = append(rightPanel, demoNotice)
	}

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
		lipgloss.NewStyle().Width(2).Render(""),
		rightContent,
	)
}

func (m Model) renderMonitorSelection(contentHeight int) string {
	var content []string

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
				detailStyle.Render(fmt.Sprintf("%s @ %.0fHz", utils.FormatResolution(monitor.Width, monitor.Height), monitor.RefreshRate)),
			)
		} else {
			card = fmt.Sprintf("  %s %s\n    %s\n    %s",
				nameStyle.Render(monitor.Name),
				statusStyle.Render(statusText),
				detailStyle.Render(fmt.Sprintf("%s %s", monitor.Make, monitor.Model)),
				detailStyle.Render(fmt.Sprintf("%s @ %.0fHz", utils.FormatResolution(monitor.Width, monitor.Height), monitor.RefreshRate)),
			)
		}

		content = append(content, card)
		content = append(content, "")
	}

	instructions := []string{
		lipgloss.NewStyle().Foreground(colorYellow).Render("⏎") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" Select monitor and return to dashboard"),
		lipgloss.NewStyle().Foreground(colorMagenta).Render("esc") +
			lipgloss.NewStyle().Foreground(colorSubtle).Render(" Return to main menu"),
	}

	note := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true).
		Render("💡 Selected monitor will be marked as CURRENT on the dashboard")

	content = append(content, "")
	content = append(content, strings.Join(instructions, "  "))
	content = append(content, "")
	content = append(content, note)

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

	title := lipgloss.NewStyle().
		Foreground(colorGreen).
		Bold(true).
		Render("🧠 Smart Scaling Recommendations")

	content = append(content, title)
	content = append(content, "")

	if len(m.monitors) > 0 && m.selectedMonitor < len(m.monitors) {
		selectedMonitor := m.monitors[m.selectedMonitor]

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

		recTitle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("🎯 Available Options")
		content = append(content, recTitle)
		content = append(content, "")

		for i, option := range m.scalingOptions {
			var line string

			if i == m.selectedScalingOpt {
				line = lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("▶ ")
			} else {
				line = "  "
			}

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

			description := fmt.Sprintf("    %s", option.Description)
			content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(description))

			details := fmt.Sprintf("    Monitor: %.1fx • GTK: %dx • Font DPI: %d • Result: %dx%d",
				option.MonitorScale, option.GTKScale, option.FontDPI,
				option.EffectiveWidth, option.EffectiveHeight)
			content = append(content, lipgloss.NewStyle().Foreground(colorComment).Render(details))

			reasoning := fmt.Sprintf("    💡 %s", option.Reasoning)
			content = append(content, lipgloss.NewStyle().Foreground(colorComment).Italic(true).Render(reasoning))

			content = append(content, "")
		}

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

	title := lipgloss.NewStyle().
		Foreground(colorMagenta).
		Bold(true).
		Render("🔧 Manual Scaling Controls")

	content = append(content, title)
	content = append(content, "")

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

	controlsTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("⚙️ Scaling Controls")
	content = append(content, controlsTitle)
	content = append(content, "")

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

	effectiveWidth, effectiveHeight := utils.CalculateEffectiveResolution(selectedMonitor.Width, selectedMonitor.Height, m.manualMonitorScale)
	screenRealEstate := utils.CalculateScreenRealEstate(m.manualMonitorScale)
	fontMultiplier := utils.CalculateFontMultiplier(m.manualFontDPI, types.BaseDPI)

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

	title := lipgloss.NewStyle().
		Foreground(colorMagenta).
		Bold(true).
		Render("⚙️ Application Settings")

	content = append(content, title)
	content = append(content, "")

	appTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("📱 Application Info")
	content = append(content, appTitle)
	content = append(content, "")

	themeInfo := m.cachedTerminalTheme
	if themeInfo == "" {
		themeInfo = getTerminalThemeInfo()
	}

	appItems := []string{
		fmt.Sprintf("  Version: %s", lipgloss.NewStyle().Foreground(colorGreen).Render("1.0.0")),
		fmt.Sprintf("  Theme: %s", lipgloss.NewStyle().Foreground(colorMagenta).Render(themeInfo)),
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

	detectionTitle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("🔍 Detection Methods")
	content = append(content, detectionTitle)
	content = append(content, "")

	commands := []string{"hyprctl", "wlr-randr"}
	names := []string{"Hyprctl", "wlr-randr"}

	for i, cmd := range commands {
		var status string

		if m.cachedCommandStatus != nil {
			if available, exists := m.cachedCommandStatus[cmd]; exists {
				if available {
					status = lipgloss.NewStyle().Foreground(colorGreen).Render("✓ Available")
				} else {
					status = lipgloss.NewStyle().Foreground(colorRed).Render("✗ Not found")
				}
			}
		}

		if status == "" {
			if utils.CommandExists(cmd) {
				status = lipgloss.NewStyle().Foreground(colorGreen).Render("✓ Available")
			} else {
				status = lipgloss.NewStyle().Foreground(colorRed).Render("✗ Not found")
			}
		}

		item := fmt.Sprintf("  %s: %s", names[i], status)
		content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(item))
	}

	content = append(content, "")

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

	footer := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true).
		Render("💡 Press Esc to return to the main menu")
	content = append(content, footer)

	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMagenta).
		Render(strings.Join(content, "\n"))
}

func (m Model) renderConfirmation(contentHeight int) string {
	var content []string

	title := lipgloss.NewStyle().
		Foreground(colorYellow).
		Bold(true).
		Render("⚠️ Confirm Scaling Changes")

	content = append(content, title)
	content = append(content, "")

	warningStyle := lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)

	warning := warningStyle.Render("⚠️  WARNING: Desktop refresh required!")
	content = append(content, warning)
	content = append(content, "")

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

	monitor := m.pendingMonitor
	option := m.pendingOption

	monitorTitle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("📱 Target Monitor")
	content = append(content, monitorTitle)
	content = append(content, "")

	monitorInfo := fmt.Sprintf("  %s (%dx%d@%.1fHz)",
		monitor.Name, monitor.Width, monitor.Height, monitor.RefreshRate)
	content = append(content, lipgloss.NewStyle().Foreground(colorSubtle).Render(monitorInfo))
	content = append(content, "")

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

	actionName := "Smart Scaling"
	if m.confirmationAction == ConfirmManualScaling {
		actionName = "Manual Scaling"
	}

	actionInfo := fmt.Sprintf("Action: %s - %s",
		lipgloss.NewStyle().Foreground(colorCyan).Render(actionName),
		lipgloss.NewStyle().Foreground(colorComment).Render(option.DisplayName))
	content = append(content, actionInfo)
	content = append(content, "")

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

	title := lipgloss.NewStyle().
		Foreground(colorYellow).
		Bold(true).
		Render("📖 Help & Controls")

	content = append(content, title)
	content = append(content, "")

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

	footer := lipgloss.NewStyle().
		Foreground(colorComment).
		Italic(true).
		Render("💡 Press Esc to return to the main menu")
	content = append(content, footer)

	return lipgloss.NewStyle().
		Width(m.width - 8).
		Height(contentHeight - 2).
		Padding(2).
		Background(colorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorYellow).
		Render(strings.Join(content, "\n"))
}
