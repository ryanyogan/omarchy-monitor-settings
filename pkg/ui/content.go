package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ContentBuilder struct {
	lines []string
}

func NewContentBuilder() *ContentBuilder {
	return &ContentBuilder{lines: make([]string, 0)}
}

func (cb *ContentBuilder) Add(line string) {
	cb.lines = append(cb.lines, line)
}

func (cb *ContentBuilder) AddEmpty() {
	cb.lines = append(cb.lines, "")
}

func (cb *ContentBuilder) AddAll(lines []string) {
	cb.lines = append(cb.lines, lines...)
}

func (cb *ContentBuilder) Build() []string {
	return cb.lines
}

func (cb *ContentBuilder) RenderJoined() string {
	return strings.Join(cb.lines, "\n")
}

type PanelRenderer struct {
	width  int
	height int
}

func NewPanelRenderer(width, height int) *PanelRenderer {
	return &PanelRenderer{width: width, height: height}
}

func (pr *PanelRenderer) RenderContent(content []string, style lipgloss.Style) string {
	return style.
		Width(pr.width).
		Height(pr.height).
		Render(strings.Join(content, "\n"))
}

func JoinInstructions(instructions []string) string {
	return strings.Join(instructions, "  ")
}

func IsCurrentlySelected(index, selectedIndex int) bool {
	return index == selectedIndex
}

func GetSelectorPrefix(isSelected bool, selector, normal string) string {
	if isSelected {
		return selector
	}
	return normal
}
