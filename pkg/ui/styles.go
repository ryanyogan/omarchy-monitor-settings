package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorComment = lipgloss.Color("8")
	colorSubtle  = lipgloss.Color("8")
	colorBlue    = lipgloss.Color("4")
	colorCyan    = lipgloss.Color("6")
	colorGreen   = lipgloss.Color("2")
	colorYellow  = lipgloss.Color("3")
	colorRed     = lipgloss.Color("1")
	colorMagenta = lipgloss.Color("5")
)

func Title(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Bold(true).Render(text)
}

func Subtitle(text string) string {
	return lipgloss.NewStyle().Foreground(colorSubtle).Render(text)
}

func Comment(text string) string {
	return lipgloss.NewStyle().Foreground(colorComment).Render(text)
}

func Colored(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Render(text)
}

func ColoredBold(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Bold(true).Render(text)
}

func Italic(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Italic(true).Render(text)
}

func KeyValue(key, value string, keyColor, valueColor lipgloss.Color) string {
	return fmt.Sprintf("  %s: %s",
		Colored(key, keyColor),
		Colored(value, valueColor))
}

func SectionTitle(emoji, title string, color lipgloss.Color) string {
	return ColoredBold(fmt.Sprintf("%s %s", emoji, title), color)
}

func NavKey(key, description string, keyColor lipgloss.Color) string {
	return fmt.Sprintf("%s %s",
		ColoredBold(key, keyColor),
		Subtitle(description))
}

func StatusAvailable() string {
	return Colored("✓ Available", colorGreen)
}

func StatusNotFound() string {
	return Colored("✗ Not found", colorRed)
}

func StatusDemo() string {
	return Colored("Demo", colorYellow)
}

func StatusLive() string {
	return Colored("Live", colorGreen)
}

func MonitorDetails(makeVal, modelVal string) string {
	return Subtitle(fmt.Sprintf("  %s %s", makeVal, modelVal))
}

func MonitorSpecs(resolution, refreshRate string) string {
	return Comment(fmt.Sprintf("  %s @ %s", resolution, refreshRate))
}

func MonitorScale(scale float64) string {
	return Comment(fmt.Sprintf("  Scale: %.1fx", scale))
}

func Selector(color lipgloss.Color) string {
	return ColoredBold("▶ ", color)
}

func ScaleValue(scale float64) string {
	return Colored(fmt.Sprintf("%.2fx", scale), colorGreen)
}

func GTKScaleValue(scale int) string {
	return Colored(fmt.Sprintf("%dx", scale), colorMagenta)
}

func DPIValue(dpi int) string {
	return Colored(fmt.Sprintf("%d", dpi), colorYellow)
}

func VersionValue(version string) string {
	return Colored(version, colorGreen)
}

func ThemeValue(theme string) string {
	return Colored(theme, colorMagenta)
}

func TargetValue(target string) string {
	return Colored(target, colorCyan)
}

func BuiltWithValue(tech string) string {
	return Colored(tech, colorBlue)
}
