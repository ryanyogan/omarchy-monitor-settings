// Package ui provides UI styling utilities and common rendering patterns.
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Colors - these are the exact colors used in model.go
var (
	colorBackground = lipgloss.Color("")
	colorSurface    = lipgloss.Color("")
	colorFloat      = lipgloss.Color("")
	colorForeground = lipgloss.Color("")
	colorComment    = lipgloss.Color("8")
	colorSubtle     = lipgloss.Color("8")
	colorBlue       = lipgloss.Color("4")
	colorCyan       = lipgloss.Color("6")
	colorGreen      = lipgloss.Color("2")
	colorYellow     = lipgloss.Color("3")
	colorRed        = lipgloss.Color("1")
	colorMagenta    = lipgloss.Color("5")
)

// Common styling patterns - EXACT SAME OUTPUT AS ORIGINAL

// Title creates a colored bold title
func Title(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Bold(true).Render(text)
}

// Subtitle creates subtle colored text
func Subtitle(text string) string {
	return lipgloss.NewStyle().Foreground(colorSubtle).Render(text)
}

// Comment creates comment-colored text
func Comment(text string) string {
	return lipgloss.NewStyle().Foreground(colorComment).Render(text)
}

// Colored creates text with a specific color
func Colored(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Render(text)
}

// ColoredBold creates bold text with a specific color
func ColoredBold(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Bold(true).Render(text)
}

// Italic creates italic text with a color
func Italic(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Italic(true).Render(text)
}

// Key-value pair styling
func KeyValue(key, value string, keyColor, valueColor lipgloss.Color) string {
	return fmt.Sprintf("  %s: %s",
		Colored(key, keyColor),
		Colored(value, valueColor))
}

// Common UI patterns with emojis and titles
func SectionTitle(emoji, title string, color lipgloss.Color) string {
	return ColoredBold(fmt.Sprintf("%s %s", emoji, title), color)
}

// Navigation key styling - exactly like original
func NavKey(key, description string, keyColor lipgloss.Color) string {
	return fmt.Sprintf("%s %s",
		ColoredBold(key, keyColor),
		Subtitle(description))
}

// Status indicators
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

// Monitor info formatting - exactly like original
func MonitorDetails(make, model string) string {
	return Subtitle(fmt.Sprintf("  %s %s", make, model))
}

func MonitorSpecs(resolution, refreshRate string) string {
	return Comment(fmt.Sprintf("  %s @ %s", resolution, refreshRate))
}

func MonitorScale(scale float64) string {
	return Comment(fmt.Sprintf("  Scale: %.1fx", scale))
}

// Selector styling
func Selector(color lipgloss.Color) string {
	return ColoredBold("▶ ", color)
}

// Common value formatting
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
