package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ContentBuilder helps build content arrays without affecting UI output
type ContentBuilder struct {
	lines []string
}

// NewContentBuilder creates a new content builder
func NewContentBuilder() *ContentBuilder {
	return &ContentBuilder{lines: make([]string, 0)}
}

// Add adds a line to the content
func (cb *ContentBuilder) Add(line string) {
	cb.lines = append(cb.lines, line)
}

// AddEmpty adds an empty line
func (cb *ContentBuilder) AddEmpty() {
	cb.lines = append(cb.lines, "")
}

// AddAll adds multiple lines
func (cb *ContentBuilder) AddAll(lines []string) {
	cb.lines = append(cb.lines, lines...)
}

// Build returns the content as a slice of strings
func (cb *ContentBuilder) Build() []string {
	return cb.lines
}

// RenderJoined renders the content joined by newlines - EXACT SAME AS strings.Join(content, "\n")
func (cb *ContentBuilder) RenderJoined() string {
	return strings.Join(cb.lines, "\n")
}

// PanelRenderer handles common panel rendering patterns
type PanelRenderer struct {
	width  int
	height int
}

// NewPanelRenderer creates a new panel renderer
func NewPanelRenderer(width, height int) *PanelRenderer {
	return &PanelRenderer{width: width, height: height}
}

// RenderContent renders content with exact same styling as original
func (pr *PanelRenderer) RenderContent(content []string, style lipgloss.Style) string {
	return style.
		Width(pr.width).
		Height(pr.height).
		Render(strings.Join(content, "\n"))
}

// InstructionFormatter formats instruction lists - EXACT SAME OUTPUT AS ORIGINAL
func JoinInstructions(instructions []string) string {
	return strings.Join(instructions, "  ")
}

// Safe utility functions that don't change output

// IsCurrentlySelected checks if an index matches the selected index
func IsCurrentlySelected(index, selectedIndex int) bool {
	return index == selectedIndex
}

// GetSelectorPrefix returns the appropriate prefix for list items
func GetSelectorPrefix(isSelected bool, selector, normal string) string {
	if isSelected {
		return selector
	}
	return normal
}
