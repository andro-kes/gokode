package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	defaultTimeout      = 5 * time.Minute
	golangciLintVersion = "v1.60.3"
	gocycloVersion      = "latest"
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
  analyse      Run full analysis (fmt, vet, lint with fixes, test, coverage, gocyclo)
  fmt          Format code with gofmt
  vet          Run go vet and write output to metrics/vet.txt
  lint         Run golangci-lint and write JSON to metrics/report.json
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
		fn   func(context.Context, string, string) int
	}{
		{"Format", func(ctx context.Context, p, m string) int { return runFormat(ctx, p) }},
		{"Vet", func(ctx context.Context, p, m string) int { return runVet(ctx, p, m) }},
		{"Lint with fixes", func(ctx context.Context, p, m string) int { return runLint(ctx, p, m, true) }},
		{"Tests", func(ctx context.Context, p, m string) int { return runTests(ctx, p) }},
		{"Coverage", func(ctx context.Context, p, m string) int { return runCoverage(ctx, p, m) }},
		{"Cyclomatic complexity", func(ctx context.Context, p, m string) int { return runGocyclo(ctx, p, m) }},
	}

	for _, step := range steps {
		fmt.Printf("\n=== %s ===\n", step.name)
		if exitCode := step.fn(ctx, path, metricsDir); exitCode != 0 {
			fmt.Fprintf(os.Stderr, "Analysis failed at step: %s\n", step.name)
			return exitCode
		}
	}

	fmt.Println("\n=== Analysis complete ===")
	fmt.Printf("Reports written to: %s\n", metricsDir)
	return 0
}

func runFormat(ctx context.Context, path string) int {
	fmt.Println("Formatting code with gofmt...")
	cmd := exec.CommandContext(ctx, "gofmt", "-w", "-s", ".")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running gofmt: %v\n", err)
		return 1
	}
	fmt.Println("✓ Format complete")
	return 0
}

func runVet(ctx context.Context, path, metricsDir string) int {
	fmt.Println("Running go vet...")
	vetFile := filepath.Join(metricsDir, "vet.txt")

	cmd := exec.CommandContext(ctx, "go", "vet", "./...")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	// Write output to file regardless of error
	if writeErr := os.WriteFile(vetFile, output, 0644); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error writing vet output: %v\n", writeErr)
		return 1
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "go vet found issues (see %s):\n%s\n", vetFile, string(output))
		// Don't fail on vet issues, just report them
	}

	fmt.Printf("✓ Vet complete (output: %s)\n", vetFile)
	return 0
}

func runLint(ctx context.Context, path, metricsDir string, fix bool) int {
	// Ensure golangci-lint is installed
	if !isToolInstalled("golangci-lint") {
		fmt.Println("golangci-lint not found, installing...")
		if err := installGolangciLint(); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing golangci-lint: %v\n", err)
			return 1
		}
	}

	reportFile := filepath.Join(metricsDir, "report.json")

	args := []string{"run", "--out-format", "json", "./..."}
	if fix {
		args = append(args, "--fix")
		fmt.Println("Running golangci-lint with --fix...")
	} else {
		fmt.Println("Running golangci-lint...")
	}

	cmd := exec.CommandContext(ctx, "golangci-lint", args...)
	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	// Write JSON output to file
	if writeErr := os.WriteFile(reportFile, output, 0644); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error writing lint report: %v\n", writeErr)
		return 1
	}

	// Also print human-readable output to console
	if len(output) > 0 && string(output) != "{}\n" && string(output) != "" {
		fmt.Println("Lint issues found:")
		// Run again without JSON for console output
		consoleCmd := exec.CommandContext(ctx, "golangci-lint", "run", "./...")
		consoleCmd.Dir = path
		consoleCmd.Stdout = os.Stdout
		consoleCmd.Stderr = os.Stderr
		_ = consoleCmd.Run() // Ignore error, we already have the JSON
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "golangci-lint found issues (see %s)\n", reportFile)
		// Don't fail on lint issues, just report them
	}

	fmt.Printf("✓ Lint complete (report: %s)\n", reportFile)
	return 0
}

func runTests(ctx context.Context, path string) int {
	fmt.Println("Running tests...")
	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-v")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Tests failed: %v\n", err)
		return 1
	}
	fmt.Println("✓ Tests passed")
	return 0
}

func runCoverage(ctx context.Context, path, metricsDir string) int {
	fmt.Println("Running tests with coverage...")
	coverageOut := filepath.Join(metricsDir, "coverage.out")
	coverageHTML := filepath.Join(metricsDir, "coverage.html")

	// Run tests with coverage
	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-coverprofile="+coverageOut)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Coverage tests failed: %v\n", err)
		return 1
	}

	// Generate HTML report
	cmd = exec.CommandContext(ctx, "go", "tool", "cover", "-html="+coverageOut, "-o", coverageHTML)
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating HTML coverage report: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Coverage complete (profile: %s, HTML: %s)\n", coverageOut, coverageHTML)
	return 0
}

func runGocyclo(ctx context.Context, path, metricsDir string) int {
	// Ensure gocyclo is installed
	if !isToolInstalled("gocyclo") {
		fmt.Println("gocyclo not found, installing...")
		if err := installGocyclo(); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing gocyclo: %v\n", err)
			return 1
		}
	}

	fmt.Println("Running cyclomatic complexity analysis...")
	gocycloFile := filepath.Join(metricsDir, "gocyclo.txt")

	cmd := exec.CommandContext(ctx, "gocyclo", ".")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	// Write output to file
	if writeErr := os.WriteFile(gocycloFile, output, 0644); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error writing gocyclo output: %v\n", writeErr)
		return 1
	}

	if err != nil {
		// gocyclo returns non-zero if it finds complex functions
		fmt.Printf("Cyclomatic complexity analysis complete (see %s)\n", gocycloFile)
	} else {
		fmt.Printf("✓ Cyclomatic complexity analysis complete (output: %s)\n", gocycloFile)
	}

	return 0
}

func installTools() int {
	fmt.Println("Installing required tools...")

	var failed bool

	if err := installGolangciLint(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing golangci-lint: %v\n", err)
		failed = true
	} else {
		fmt.Println("✓ golangci-lint installed")
	}

	if err := installGocyclo(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing gocyclo: %v\n", err)
		failed = true
	} else {
		fmt.Println("✓ gocyclo installed")
	}

	if failed {
		return 1
	}

	fmt.Println("✓ All tools installed successfully")
	return 0
}

func isToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

func installGolangciLint() error {
	fmt.Printf("Installing golangci-lint %s...\n", golangciLintVersion)
	cmd := exec.Command("go", "install", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", golangciLintVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installGocyclo() error {
	fmt.Printf("Installing gocyclo %s...\n", gocycloVersion)
	cmd := exec.Command("go", "install", "github.com/fzipp/gocyclo/cmd/gocyclo@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
