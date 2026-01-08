package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/andro-kes/gokode/internal/tools"
)

// RunFormat formats code with gofmt
func RunFormat(ctx context.Context, path string) error {
	fmt.Println("Formatting code with gofmt...")
	cmd := exec.CommandContext(ctx, "gofmt", "-w", "-s", ".")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running gofmt: %w", err)
	}
	fmt.Println("✓ Format complete")
	return nil
}

// RunVet runs go vet and writes output to a file
func RunVet(ctx context.Context, path, metricsDir string) error {
	fmt.Println("Running go vet...")
	vetFile := filepath.Join(metricsDir, "vet.txt")

	cmd := exec.CommandContext(ctx, "go", "vet", "./...")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	// Write output to file regardless of error
	if writeErr := os.WriteFile(vetFile, output, 0644); writeErr != nil {
		return fmt.Errorf("error writing vet output: %w", writeErr)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "go vet found issues (see %s):\n%s\n", vetFile, string(output))
		// Don't fail on vet issues, just report them
	}

	fmt.Printf("✓ Vet complete (output: %s)\n", vetFile)
	return nil
}

// RunLint runs golangci-lint and writes pretty-printed JSON to a file
func RunLint(ctx context.Context, path, metricsDir string, fix bool) error {
	// Ensure golangci-lint is installed
	if !tools.IsInstalled("golangci-lint") {
		fmt.Println("golangci-lint not found, installing...")
		if err := tools.InstallGolangciLint(); err != nil {
			return fmt.Errorf("error installing golangci-lint: %w", err)
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

	// Parse and pretty-print the JSON output
	var jsonData interface{}
	if len(output) > 0 {
		if jsonErr := json.Unmarshal(output, &jsonData); jsonErr == nil {
			// Successfully parsed JSON, now pretty-print it
			prettyJSON, marshalErr := json.MarshalIndent(jsonData, "", "  ")
			if marshalErr == nil {
				output = prettyJSON
			}
		}
	}

	// Write JSON output to file
	if writeErr := os.WriteFile(reportFile, output, 0644); writeErr != nil {
		return fmt.Errorf("error writing lint report: %w", writeErr)
	}

	// Also print human-readable output to console
	if len(output) > 0 {
		// Check if JSON has any actual issues
		var lintReport struct {
			Issues []interface{} `json:"Issues"`
		}
		hasIssues := false
		if json.Unmarshal(output, &lintReport) == nil && len(lintReport.Issues) > 0 {
			hasIssues = true
		}

		if hasIssues {
			fmt.Println("Lint issues found:")
			// Run again without JSON for console output
			consoleCmd := exec.CommandContext(ctx, "golangci-lint", "run", "./...")
			consoleCmd.Dir = path
			consoleCmd.Stdout = os.Stdout
			consoleCmd.Stderr = os.Stderr
			_ = consoleCmd.Run() // Ignore error, we already have the JSON
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "golangci-lint found issues (see %s)\n", reportFile)
		// Don't fail on lint issues, just report them
	}

	fmt.Printf("✓ Lint complete (report: %s)\n", reportFile)
	return nil
}

// RunTests runs go tests
func RunTests(ctx context.Context, path string) error {
	fmt.Println("Running tests...")
	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-v")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}
	fmt.Println("✓ Tests passed")
	return nil
}

// RunCoverage runs tests with coverage and generates HTML report
func RunCoverage(ctx context.Context, path, metricsDir string) error {
	fmt.Println("Running tests with coverage...")
	coverageOut := filepath.Join(metricsDir, "coverage.out")
	coverageHTML := filepath.Join(metricsDir, "coverage.html")

	// Run tests with coverage
	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-coverprofile="+coverageOut)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("coverage tests failed: %w", err)
	}

	// Generate HTML report
	cmd = exec.CommandContext(ctx, "go", "tool", "cover", "-html="+coverageOut, "-o", coverageHTML)
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error generating HTML coverage report: %w", err)
	}

	fmt.Printf("✓ Coverage complete (profile: %s, HTML: %s)\n", coverageOut, coverageHTML)
	return nil
}

// RunGocyclo runs cyclomatic complexity analysis
func RunGocyclo(ctx context.Context, path, metricsDir string) error {
	// Ensure gocyclo is installed
	if !tools.IsInstalled("gocyclo") {
		fmt.Println("gocyclo not found, installing...")
		if err := tools.InstallGocyclo(); err != nil {
			return fmt.Errorf("error installing gocyclo: %w", err)
		}
	}

	fmt.Println("Running cyclomatic complexity analysis...")
	gocycloFile := filepath.Join(metricsDir, "gocyclo.txt")

	cmd := exec.CommandContext(ctx, "gocyclo", ".")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	// Write output to file
	if writeErr := os.WriteFile(gocycloFile, output, 0644); writeErr != nil {
		return fmt.Errorf("error writing gocyclo output: %w", writeErr)
	}

	if err != nil {
		// gocyclo returns non-zero if it finds complex functions
		fmt.Printf("Cyclomatic complexity analysis complete (see %s)\n", gocycloFile)
	} else {
		fmt.Printf("✓ Cyclomatic complexity analysis complete (output: %s)\n", gocycloFile)
	}

	return nil
}
