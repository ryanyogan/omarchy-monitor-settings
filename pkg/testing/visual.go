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

type VisualTestConfig struct {
	Name         string
	Width        int
	Height       int
	Model        tea.Model
	UpdateGolden bool
}

type VisualSnapshot struct {
	Name     string
	Width    int
	Height   int
	Content  string
	Hash     string
	Metadata map[string]interface{}
}

type VisualTester struct {
	goldenDir string
	t         *testing.T
}

func NewVisualTester(t *testing.T, goldenDir string) *VisualTester {
	if goldenDir == "" {
		goldenDir = "testdata/golden"
	}

	if err := os.MkdirAll(goldenDir, 0750); err != nil {
		t.Fatalf("Failed to create golden directory: %v", err)
	}

	return &VisualTester{
		goldenDir: goldenDir,
		t:         t,
	}
}

func (vt *VisualTester) CaptureSnapshot(config VisualTestConfig) *VisualSnapshot {
	model := config.Model

	model, _ = model.Update(tea.WindowSizeMsg{
		Width:  config.Width,
		Height: config.Height,
	})

	content := model.View()

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

func (vt *VisualTester) CompareWithGolden(snapshot *VisualSnapshot, updateGolden bool) error {
	goldenPath := filepath.Join(vt.goldenDir, fmt.Sprintf("%s_%dx%d.golden",
		snapshot.Name, snapshot.Width, snapshot.Height))

	if updateGolden {
		return vt.updateGoldenFile(goldenPath, snapshot)
	}

	return vt.compareWithGoldenFile(goldenPath, snapshot)
}

func (vt *VisualTester) updateGoldenFile(goldenPath string, snapshot *VisualSnapshot) error {
	content := fmt.Sprintf("# Visual Golden File\n# Name: %s\n# Dimensions: %dx%d\n# Hash: %s\n\n%s",
		snapshot.Name, snapshot.Width, snapshot.Height, snapshot.Hash, snapshot.Content)

	return os.WriteFile(goldenPath, []byte(content), 0600)
}

func (vt *VisualTester) compareWithGoldenFile(goldenPath string, snapshot *VisualSnapshot) error {
	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		return fmt.Errorf("golden file does not exist: %s\nRun with UPDATE_GOLDEN=true to create it", goldenPath)
	}

	if !strings.HasPrefix(goldenPath, vt.goldenDir) {
		return fmt.Errorf("golden file path is outside allowed directory")
	}
	goldenBytes, err := os.ReadFile(goldenPath) // nosec G304
	if err != nil {
		return fmt.Errorf("failed to read golden file: %v", err)
	}

	goldenContent := string(goldenBytes)

	parts := strings.SplitN(goldenContent, "\n\n", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid golden file format")
	}

	expectedContent := parts[1]

	if snapshot.Content != expectedContent {
		diffPath := filepath.Join(vt.goldenDir, fmt.Sprintf("%s_%dx%d.diff",
			snapshot.Name, snapshot.Width, snapshot.Height))

		diffContent := fmt.Sprintf("Visual regression detected!\n\n=== EXPECTED ===\n%s\n\n=== ACTUAL ===\n%s\n\n=== DIFF INFO ===\nExpected length: %d\nActual length: %d\nExpected hash: %s\nActual hash: %s\n",
			expectedContent, snapshot.Content, len(expectedContent), len(snapshot.Content),
			vt.calculateHash(expectedContent), snapshot.Hash)

		if err := os.WriteFile(diffPath, []byte(diffContent), 0600); err != nil {
			// Log error but don't fail the test
			fmt.Printf("Warning: failed to write diff file: %v\n", err)
		}

		return fmt.Errorf("visual regression detected in %s\nExpected hash: %s\nActual hash: %s\nDiff saved to: %s",
			snapshot.Name, vt.calculateHash(expectedContent), snapshot.Hash, diffPath)
	}

	return nil
}

func (vt *VisualTester) calculateHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

func (vt *VisualTester) TestVisualRegression(config VisualTestConfig) {
	snapshot := vt.CaptureSnapshot(config)

	updateGolden := config.UpdateGolden || os.Getenv("UPDATE_GOLDEN") == "true"

	if err := vt.CompareWithGolden(snapshot, updateGolden); err != nil {
		vt.t.Errorf("Visual regression test failed: %v", err)
	}
}

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

func StripANSI(input string) string {
	return lipgloss.NewStyle().Render(input)
}

func NormalizeWhitespace(input string) string {
	lines := strings.Split(input, "\n")
	var normalized []string

	for _, line := range lines {
		normalized = append(normalized, strings.TrimRight(line, " \t"))
	}

	return strings.Join(normalized, "\n")
}
