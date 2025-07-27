# Visual Regression Testing

This document explains the visual regression testing framework designed to ensure UI consistency during code refactoring and changes.

## Overview

Visual regression testing captures the terminal output of the application in various states and compares it against known-good "golden" files. This prevents unintended UI changes during refactoring or feature development.

## Framework Components

### 1. Visual Testing Package (`pkg/testing/visual.go`)

The core framework provides:

- **VisualTester**: Main testing orchestrator
- **VisualSnapshot**: Captures UI state with metadata
- **Golden Files**: Reference outputs stored in `testdata/golden/`
- **Diff Generation**: Detailed comparison reports

### 2. Test Suite (`visual_regression_test.go`)

Comprehensive tests covering:

- **All UI States**: Dashboard, Monitor Selection, Scaling Options, Manual Scaling, Settings, Help
- **Multiple Screen Sizes**: 80x24 (minimum) to 200x60 (large terminal)
- **Edge Cases**: Empty monitors, terminal too small, many monitors
- **Interactive States**: Different menu selections, scaling values
- **Terminal Themes**: xterm, screen, tmux compatibility

## Usage

### Running Visual Tests

```bash
# Run all visual regression tests
make visual-test

# Run specific test categories
go test -run TestVisualRegression/Dashboard -v
go test -run TestVisualRegressionEdgeCases -v
go test -run TestVisualRegressionInteractions -v
```

### Updating Golden Files

When UI changes are intentional:

```bash
# Update all golden files
make visual-update

# Update specific golden files
UPDATE_GOLDEN=true go test -run TestVisualRegression/Dashboard -v
```

### Cleaning Test Artifacts

```bash
# Remove diff files and temporary artifacts
make visual-clean
```

## Golden File Structure

Golden files are stored in `testdata/golden/` with the naming convention:
```
{test_name}_{width}x{height}.golden
```

Example files:
- `dashboard_120x40.golden` - Dashboard view at 120x40 resolution
- `manual_scaling_control_1_120x40.golden` - Manual scaling with GTK scale selected
- `terminal_too_small_70x20.golden` - Terminal below minimum size

## File Format

```
# Visual Golden File
# Name: dashboard
# Dimensions: 120x40
# Hash: a1b2c3d4e5f6...

[Actual terminal output content...]
```

## Testing Strategy

### 1. Core UI States
- Tests all major application modes
- Ensures consistent layout across screen sizes
- Validates responsive behavior

### 2. Edge Cases
- Empty data states (no monitors)
- Boundary conditions (terminal too small)
- Large datasets (many monitors)

### 3. Interactive States
- Menu navigation selections
- Control focus states
- Different value ranges

### 4. Terminal Compatibility
- Different TERM environment variables
- Color support variations
- Terminal-specific behaviors

## Best Practices

### 1. When to Update Golden Files
- ✅ **Intentional UI improvements**
- ✅ **New features that change layout**
- ✅ **Fixes to UI bugs**
- ❌ **Accidental changes during refactoring**

### 2. Reviewing Changes
Before updating golden files:
1. Run `make visual-test` to see what changed
2. Check generated `.diff` files in `testdata/golden/`
3. Verify changes are intentional
4. Update golden files only if changes are correct

### 3. CI/CD Integration
```bash
# In CI pipeline
make visual-test  # Fail if visual regressions detected
```

## Debugging Visual Failures

When visual regression tests fail:

1. **Check the diff file**: `testdata/golden/{test_name}.diff`
2. **Compare expected vs actual output**
3. **Look for**:
   - Changed colors or styling
   - Layout shifts
   - Text content changes
   - Spacing modifications

### Common Failure Causes

- **Color changes**: ANSI color code modifications
- **Spacing**: Padding, margin, or width changes
- **Content**: Text modifications or additions
- **Layout**: Panel size or positioning changes

## Example Workflow

### Adding a New Feature
1. Implement feature
2. Run `make visual-test` - expect failures
3. Review `.diff` files to ensure changes are correct
4. Run `make visual-update` to accept changes
5. Commit both code and updated golden files

### Refactoring Code
1. Make refactoring changes
2. Run `make visual-test` - should pass
3. If failures occur, investigate:
   - Are they unintended UI changes?
   - Fix the code to maintain UI consistency

### Debugging Test Failures
1. Run specific failing test: `go test -run TestVisualRegression/Dashboard -v`
2. Check generated diff: `cat testdata/golden/dashboard_120x40.diff`
3. Compare expected vs actual output
4. Fix code or update golden file as appropriate

## Framework Benefits

1. **Prevents UI Regressions**: Catches unintended visual changes
2. **Enables Safe Refactoring**: Confidence that UI remains unchanged
3. **Multiple Screen Size Testing**: Ensures responsive behavior
4. **Comprehensive Coverage**: Tests all application states
5. **Easy Debugging**: Clear diff output shows exactly what changed
6. **CI/CD Ready**: Automated testing in build pipelines

## Technical Details

### Snapshot Capture
- Renders complete TUI model at specified dimensions
- Captures raw terminal output including ANSI codes
- Calculates SHA256 hash for quick comparison

### Comparison Algorithm
- Byte-for-byte comparison of terminal output
- Preserves all formatting, colors, and spacing
- Generates detailed diffs when mismatches occur

### Mock Services
- Uses mock implementations for consistent test data
- Isolated from real system dependencies
- Predictable output across test runs

This visual regression testing framework ensures that your beautiful TUI remains pixel-perfect during development! 🎨✨ 