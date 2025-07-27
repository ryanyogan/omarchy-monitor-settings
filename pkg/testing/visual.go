// Package testing provides visual regression testing utilities for terminal UIs.
package testing

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VisualTestConfig holds configuration for visual regression tests.
type VisualTestConfig struct {
	Name         string
	Width        int
	Height       int
	Model        tea.Model
	UpdateGolden bool // Set to true to update golden files
}

// VisualSnapshot represents a captured UI state.
type VisualSnapshot struct {
	Name     string
	Width    int
	Height   int
	Content  string
	Hash     string
	Metadata map[string]interface{}
}

// VisualTester manages visual regression testing.
type VisualTester struct {
	goldenDir string
	t         *testing.T
}

// NewVisualTester creates a new visual testing instance.
func NewVisualTester(t *testing.T, goldenDir string) *VisualTester {
	if goldenDir == "" {
		goldenDir = "testdata/golden"
	}

	// Ensure golden directory exists
	if err := os.MkdirAll(goldenDir, 0755); err != nil {
		t.Fatalf("Failed to create golden directory: %v", err)
	}

	return &VisualTester{
		goldenDir: goldenDir,
		t:         t,
	}
}

// CaptureSnapshot captures the current visual state of a model.
func (vt *VisualTester) CaptureSnapshot(config VisualTestConfig) *VisualSnapshot {
	// Set up the model with the specified dimensions
	model := config.Model

	// Send window size message to properly initialize the model
	model, _ = model.Update(tea.WindowSizeMsg{
		Width:  config.Width,
		Height: config.Height,
	})

	// Capture the rendered output
	content := model.View()

	// Calculate hash for quick comparison
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	return &VisualSnapshot{
		Name:    config.Name,
		Width:   config.Width,
		Height:  config.Height,
		Content: content,
		Hash:    hash,
		Metadata: map[string]interface{}{
			"terminal_width":  config.Width,
			"terminal_height": config.Height,
			"content_length":  len(content),
			"line_count":      strings.Count(content, "\n") + 1,
		},
	}
}

// CompareWithGolden compares a snapshot with the golden file.
func (vt *VisualTester) CompareWithGolden(snapshot *VisualSnapshot, updateGolden bool) error {
	goldenPath := filepath.Join(vt.goldenDir, fmt.Sprintf("%s_%dx%d.golden",
		snapshot.Name, snapshot.Width, snapshot.Height))

	if updateGolden {
		return vt.updateGoldenFile(goldenPath, snapshot)
	}

	return vt.compareWithGoldenFile(goldenPath, snapshot)
}

// updateGoldenFile updates or creates a golden file.
func (vt *VisualTester) updateGoldenFile(goldenPath string, snapshot *VisualSnapshot) error {
	content := fmt.Sprintf("# Visual Golden File\n# Name: %s\n# Dimensions: %dx%d\n# Hash: %s\n\n%s",
		snapshot.Name, snapshot.Width, snapshot.Height, snapshot.Hash, snapshot.Content)

	return os.WriteFile(goldenPath, []byte(content), 0644)
}

// compareWithGoldenFile compares snapshot with existing golden file.
func (vt *VisualTester) compareWithGoldenFile(goldenPath string, snapshot *VisualSnapshot) error {
	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		return fmt.Errorf("golden file does not exist: %s\nRun with UPDATE_GOLDEN=true to create it", goldenPath)
	}

	goldenBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		return fmt.Errorf("failed to read golden file: %v", err)
	}

	goldenContent := string(goldenBytes)

	// Extract the actual content (everything after the double newline)
	parts := strings.SplitN(goldenContent, "\n\n", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid golden file format")
	}

	expectedContent := parts[1]

	if snapshot.Content != expectedContent {
		// Create a detailed diff for debugging
		diffPath := filepath.Join(vt.goldenDir, fmt.Sprintf("%s_%dx%d.diff",
			snapshot.Name, snapshot.Width, snapshot.Height))

		diffContent := fmt.Sprintf("Visual regression detected!\n\n=== EXPECTED ===\n%s\n\n=== ACTUAL ===\n%s\n\n=== DIFF INFO ===\nExpected length: %d\nActual length: %d\nExpected hash: %s\nActual hash: %s\n",
			expectedContent, snapshot.Content, len(expectedContent), len(snapshot.Content),
			vt.calculateHash(expectedContent), snapshot.Hash)

		os.WriteFile(diffPath, []byte(diffContent), 0644)

		return fmt.Errorf("visual regression detected in %s\nExpected hash: %s\nActual hash: %s\nDiff saved to: %s",
			snapshot.Name, vt.calculateHash(expectedContent), snapshot.Hash, diffPath)
	}

	return nil
}

// calculateHash calculates SHA256 hash of content.
func (vt *VisualTester) calculateHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

// TestVisualRegression runs a visual regression test.
func (vt *VisualTester) TestVisualRegression(config VisualTestConfig) {
	snapshot := vt.CaptureSnapshot(config)

	updateGolden := config.UpdateGolden || os.Getenv("UPDATE_GOLDEN") == "true"

	if err := vt.CompareWithGolden(snapshot, updateGolden); err != nil {
		vt.t.Errorf("Visual regression test failed: %v", err)
	}
}

// MultiSizeTest tests multiple screen sizes for responsive behavior.
func (vt *VisualTester) MultiSizeTest(name string, model tea.Model, sizes []struct{ Width, Height int }) {
	for _, size := range sizes {
		testName := fmt.Sprintf("%s_%dx%d", name, size.Width, size.Height)
		config := VisualTestConfig{
			Name:   testName,
			Width:  size.Width,
			Height: size.Height,
			Model:  model,
		}
		vt.TestVisualRegression(config)
	}
}

// StripANSI removes ANSI escape codes for content-only comparison.
func StripANSI(input string) string {
	// Use lipgloss's built-in ANSI stripping
	return lipgloss.NewStyle().Render(input)
}

// NormalizeWhitespace normalizes whitespace for more stable comparisons.
func NormalizeWhitespace(input string) string {
	lines := strings.Split(input, "\n")
	var normalized []string

	for _, line := range lines {
		// Trim trailing whitespace but preserve leading whitespace (important for TUI layout)
		normalized = append(normalized, strings.TrimRight(line, " \t"))
	}

	return strings.Join(normalized, "\n")
}
