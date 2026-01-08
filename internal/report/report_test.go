package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateHTML(t *testing.T) {
	// Create a temporary directory for test metrics
	tmpDir := t.TempDir()
	metricsDir := filepath.Join(tmpDir, "metrics")
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		t.Fatalf("Failed to create metrics dir: %v", err)
	}

	// Create sample metric files
	vetContent := "# command-line-arguments\n./main.go:10:2: unused variable x\n"
	if err := os.WriteFile(filepath.Join(metricsDir, "vet.txt"), []byte(vetContent), 0644); err != nil {
		t.Fatalf("Failed to write vet.txt: %v", err)
	}

	lintJSON := `{
  "Issues": [
    {
      "FromLinter": "errcheck",
      "Text": "Error return value not checked",
      "Pos": {
        "Filename": "main.go",
        "Line": 15,
        "Column": 5
      }
    }
  ]
}`
	if err := os.WriteFile(filepath.Join(metricsDir, "report.json"), []byte(lintJSON), 0644); err != nil {
		t.Fatalf("Failed to write report.json: %v", err)
	}

	coverageContent := "mode: set\ngithub.com/example/pkg/file.go:10.1,12.2 1 1\n"
	if err := os.WriteFile(filepath.Join(metricsDir, "coverage.out"), []byte(coverageContent), 0644); err != nil {
		t.Fatalf("Failed to write coverage.out: %v", err)
	}

	gocycloContent := "10 main main.go:15:1\n5 helper utils.go:20:1\n"
	if err := os.WriteFile(filepath.Join(metricsDir, "gocyclo.txt"), []byte(gocycloContent), 0644); err != nil {
		t.Fatalf("Failed to write gocyclo.txt: %v", err)
	}

	// Generate HTML report
	if err := GenerateHTML(metricsDir); err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Verify report.html was created
	reportPath := filepath.Join(metricsDir, "report.html")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Fatal("report.html was not created")
	}

	// Read and verify HTML content
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report.html: %v", err)
	}

	htmlString := string(content)

	// Check for key content
	expectedStrings := []string{
		"<!DOCTYPE html>",
		"Отчет анализа кода gokode",
		"Go Vet",
		"Golangci-lint",
		"Покрытие тестами",
		"Цикломатическая сложность",
		"errcheck",
		"Error return value not checked",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(htmlString, expected) {
			t.Errorf("HTML report missing expected content: %s", expected)
		}
	}
}

func TestGenerateHTMLEmptyMetrics(t *testing.T) {
	// Create a temporary directory with empty metrics
	tmpDir := t.TempDir()
	metricsDir := filepath.Join(tmpDir, "metrics")
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		t.Fatalf("Failed to create metrics dir: %v", err)
	}

	// Create empty metric files
	emptyFiles := []string{"vet.txt", "report.json", "gocyclo.txt"}
	for _, file := range emptyFiles {
		if err := os.WriteFile(filepath.Join(metricsDir, file), []byte(""), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", file, err)
		}
	}

	// Generate HTML report
	if err := GenerateHTML(metricsDir); err != nil {
		t.Fatalf("GenerateHTML failed with empty metrics: %v", err)
	}

	// Verify report.html was created
	reportPath := filepath.Join(metricsDir, "report.html")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Fatal("report.html was not created for empty metrics")
	}
}
