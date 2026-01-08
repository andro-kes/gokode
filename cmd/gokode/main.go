package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andro-kes/gokode/internal/report"
	"github.com/andro-kes/gokode/internal/runner"
	"github.com/andro-kes/gokode/internal/tools"
)

const (
	defaultTimeout = 5 * time.Minute
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	path := "."
	if len(os.Args) > 2 {
		path = os.Args[2]
	}

	// Ensure path exists
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path %s: %v\n", path, err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: path does not exist: %s\n", absPath)
		os.Exit(1)
	}

	// Create metrics directory
	metricsDir := filepath.Join(absPath, "metrics")
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating metrics directory: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var exitCode int
	switch command {
	case "analyse":
		exitCode = runAnalyse(ctx, absPath, metricsDir)
	case "fmt":
		exitCode = runFormat(ctx, absPath)
	case "vet":
		exitCode = runVet(ctx, absPath, metricsDir)
	case "lint":
		exitCode = runLint(ctx, absPath, metricsDir, false)
	case "lint-fix":
		exitCode = runLint(ctx, absPath, metricsDir, true)
	case "test":
		exitCode = runTests(ctx, absPath)
	case "coverage":
		exitCode = runCoverage(ctx, absPath, metricsDir)
	case "gocyclo":
		exitCode = runGocyclo(ctx, absPath, metricsDir)
	case "tools":
		exitCode = installTools()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		exitCode = 1
	}

	os.Exit(exitCode)
}

func printUsage() {
	usage := `gokode - Go code analysis and quality tool

Usage:
  gokode <command> [path]

Commands:
  analyse      Run full analysis (fmt, vet, lint with fixes, test, coverage, gocyclo) and generate HTML report
  fmt          Format code with gofmt
  vet          Run go vet and write output to metrics/vet.txt
  lint         Run golangci-lint and write pretty-printed JSON to metrics/report.json
  lint-fix     Run golangci-lint with --fix
  test         Run tests
  coverage     Run tests with coverage (metrics/coverage.out and coverage.html)
  gocyclo      Run cyclomatic complexity analysis (metrics/gocyclo.txt)
  tools        Install required tools (golangci-lint, gocyclo)

Arguments:
  path         Target directory (default: current directory)

Examples:
  gokode analyse .
  gokode lint ./myproject
  gokode coverage /path/to/project
`
	fmt.Fprint(os.Stderr, usage)
}

func runAnalyse(ctx context.Context, path, metricsDir string) int {
	fmt.Println("Starting full analysis...")

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Format", func() error { return runner.RunFormat(ctx, path) }},
		{"Vet", func() error { return runner.RunVet(ctx, path, metricsDir) }},
		{"Lint with fixes", func() error { return runner.RunLint(ctx, path, metricsDir, true) }},
		{"Tests", func() error { return runner.RunTests(ctx, path) }},
		{"Coverage", func() error { return runner.RunCoverage(ctx, path, metricsDir) }},
		{"Cyclomatic complexity", func() error { return runner.RunGocyclo(ctx, path, metricsDir) }},
	}

	for _, step := range steps {
		fmt.Printf("\n=== %s ===\n", step.name)
		if err := step.fn(); err != nil {
			fmt.Fprintf(os.Stderr, "Analysis failed at step: %s: %v\n", step.name, err)
			return 1
		}
	}

	// Generate HTML report
	fmt.Println("\n=== Generating HTML report ===")
	if err := report.GenerateHTML(metricsDir); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to generate HTML report: %v\n", err)
		// Don't fail the entire analysis if HTML generation fails
	}

	fmt.Println("\n=== Analysis complete ===")
	fmt.Printf("Reports written to: %s\n", metricsDir)
	return 0
}

func runFormat(ctx context.Context, path string) int {
	if err := runner.RunFormat(ctx, path); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func runVet(ctx context.Context, path, metricsDir string) int {
	if err := runner.RunVet(ctx, path, metricsDir); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func runLint(ctx context.Context, path, metricsDir string, fix bool) int {
	if err := runner.RunLint(ctx, path, metricsDir, fix); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func runTests(ctx context.Context, path string) int {
	if err := runner.RunTests(ctx, path); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func runCoverage(ctx context.Context, path, metricsDir string) int {
	if err := runner.RunCoverage(ctx, path, metricsDir); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func runGocyclo(ctx context.Context, path, metricsDir string) int {
	if err := runner.RunGocyclo(ctx, path, metricsDir); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func installTools() int {
	if err := tools.InstallAll(); err != nil {
		return 1
	}
	return 0
}
