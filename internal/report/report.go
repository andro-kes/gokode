package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MetricsSummary contains aggregated metrics data
type MetricsSummary struct {
	Timestamp      string
	VetOutput      string
	VetIssueCount  int
	LintIssues     []LintIssue
	LintIssueCount int
	CoverageData   string
	CoverageHTML   string
	GocycloOutput  string
	GocycloLines   []string
	MetricsDir     string
}

// LintIssue represents a single linting issue
type LintIssue struct {
	FromLinter  string   `json:"fromLinter"`
	Text        string   `json:"text"`
	SourceLines []string `json:"sourceLines"`
	Pos         struct {
		Filename string `json:"filename"`
		Line     int    `json:"line"`
		Column   int    `json:"column"`
	} `json:"pos"`
}

// LintReport represents the golangci-lint JSON report structure
type LintReport struct {
	Issues []LintIssue `json:"Issues"`
}

// GenerateHTML generates an HTML report from metrics files
func GenerateHTML(metricsDir string) error {
	fmt.Println("Generating HTML report...")

	summary, err := collectMetrics(metricsDir)
	if err != nil {
		return fmt.Errorf("error collecting metrics: %w", err)
	}

	htmlContent, err := renderHTML(summary)
	if err != nil {
		return fmt.Errorf("error rendering HTML: %w", err)
	}

	reportPath := filepath.Join(metricsDir, "report.html")
	if err := os.WriteFile(reportPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("error writing HTML report: %w", err)
	}

	fmt.Printf("✓ HTML report generated: %s\n", reportPath)
	return nil
}

func collectMetrics(metricsDir string) (*MetricsSummary, error) {
	summary := &MetricsSummary{
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		MetricsDir: metricsDir,
	}

	// Read vet output
	vetFile := filepath.Join(metricsDir, "vet.txt")
	if data, err := os.ReadFile(vetFile); err == nil {
		summary.VetOutput = string(data)
		trimmed := strings.TrimSpace(summary.VetOutput)
		if trimmed != "" {
			// Count non-empty lines as issues
			lines := strings.Split(trimmed, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					summary.VetIssueCount++
				}
			}
		}
	}

	// Read lint report
	lintFile := filepath.Join(metricsDir, "report.json")
	if data, err := os.ReadFile(lintFile); err == nil && len(data) > 0 {
		var lintReport LintReport
		if json.Unmarshal(data, &lintReport) == nil {
			summary.LintIssues = lintReport.Issues
			summary.LintIssueCount = len(lintReport.Issues)
		}
	}

	// Read coverage data
	coverageFile := filepath.Join(metricsDir, "coverage.out")
	if data, err := os.ReadFile(coverageFile); err == nil {
		summary.CoverageData = string(data)
	}

	// Check if coverage HTML exists
	coverageHTML := filepath.Join(metricsDir, "coverage.html")
	if _, err := os.Stat(coverageHTML); err == nil {
		summary.CoverageHTML = "coverage.html"
	}

	// Read gocyclo output
	gocycloFile := filepath.Join(metricsDir, "gocyclo.txt")
	if data, err := os.ReadFile(gocycloFile); err == nil {
		summary.GocycloOutput = string(data)
		trimmed := strings.TrimSpace(summary.GocycloOutput)
		if trimmed != "" {
			summary.GocycloLines = strings.Split(trimmed, "\n")
		}
	}

	return summary, nil
}

func renderHTML(summary *MetricsSummary) (string, error) {
	tmpl := template.Must(template.New("report").Parse(htmlTemplate))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, summary); err != nil {
		return "", err
	}

	return buf.String(), nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Отчет анализа кода gokode</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        .timestamp {
            opacity: 0.9;
            font-size: 0.95em;
        }
        .content {
            padding: 30px;
        }
        .section {
            margin-bottom: 40px;
        }
        .section h2 {
            color: #667eea;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #667eea;
            font-size: 1.8em;
        }
        .metric-card {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 20px;
            margin-bottom: 20px;
            border-radius: 4px;
        }
        .metric-card h3 {
            color: #495057;
            margin-bottom: 10px;
            font-size: 1.3em;
        }
        .status-ok {
            color: #28a745;
            font-weight: bold;
        }
        .status-warning {
            color: #ffc107;
            font-weight: bold;
        }
        .status-error {
            color: #dc3545;
            font-weight: bold;
        }
        .issue-count {
            display: inline-block;
            background: #dc3545;
            color: white;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.9em;
            margin-left: 10px;
        }
        .issue-count.ok {
            background: #28a745;
        }
        pre {
            background: #f4f4f4;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 15px;
            overflow-x: auto;
            font-size: 0.9em;
            line-height: 1.4;
        }
        .issue-item {
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 4px;
            padding: 15px;
            margin-bottom: 10px;
        }
        .issue-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .issue-linter {
            background: #667eea;
            color: white;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 0.85em;
        }
        .issue-location {
            color: #6c757d;
            font-size: 0.9em;
        }
        .issue-text {
            color: #495057;
            margin-bottom: 10px;
        }
        .issue-code {
            background: #f8f9fa;
            border-left: 3px solid #dc3545;
            padding: 10px;
            font-family: 'Courier New', monospace;
            font-size: 0.85em;
            overflow-x: auto;
        }
        .links {
            margin-top: 20px;
        }
        .links a {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 4px;
            margin-right: 10px;
            margin-bottom: 10px;
            transition: background 0.3s;
        }
        .links a:hover {
            background: #764ba2;
        }
        footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #6c757d;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 Отчет анализа кода gokode</h1>
            <div class="timestamp">Сгенерирован: {{.Timestamp}}</div>
        </header>

        <div class="content">
            <!-- Vet Section -->
            <div class="section">
                <h2>🔍 Go Vet</h2>
                <div class="metric-card">
                    <h3>Статус: {{if eq .VetIssueCount 0}}<span class="status-ok">✓ Проблем не обнаружено</span>{{else}}<span class="status-error">✗ Обнаружено проблем</span><span class="issue-count">{{.VetIssueCount}}</span>{{end}}</h3>
                    {{if .VetOutput}}
                    <pre>{{.VetOutput}}</pre>
                    {{else}}
                    <p class="status-ok">Анализ go vet завершился успешно, проблем не найдено.</p>
                    {{end}}
                </div>
            </div>

            <!-- Lint Section -->
            <div class="section">
                <h2>🔎 Golangci-lint</h2>
                <div class="metric-card">
                    <h3>Статус: {{if eq .LintIssueCount 0}}<span class="status-ok">✓ Проблем не обнаружено</span>{{else}}<span class="status-error">✗ Обнаружено проблем</span><span class="issue-count">{{.LintIssueCount}}</span>{{end}}</h3>
                    {{if .LintIssues}}
                    {{range .LintIssues}}
                    <div class="issue-item">
                        <div class="issue-header">
                            <span class="issue-linter">{{.FromLinter}}</span>
                            <span class="issue-location">{{.Pos.Filename}}:{{.Pos.Line}}:{{.Pos.Column}}</span>
                        </div>
                        <div class="issue-text">{{.Text}}</div>
                        {{if .SourceLines}}
                        <div class="issue-code">{{range .SourceLines}}{{.}}{{"\n"}}{{end}}</div>
                        {{end}}
                    </div>
                    {{end}}
                    {{else}}
                    <p class="status-ok">Анализ golangci-lint завершился успешно, проблем не найдено.</p>
                    {{end}}
                </div>
                <div class="links">
                    <a href="report.json" target="_blank">📄 Смотреть JSON отчет</a>
                </div>
            </div>

            <!-- Coverage Section -->
            <div class="section">
                <h2>📈 Покрытие тестами</h2>
                <div class="metric-card">
                    <h3>Статус: {{if .CoverageData}}<span class="status-ok">✓ Отчет о покрытии сгенерирован</span>{{else}}<span class="status-warning">⚠ Данные о покрытии отсутствуют</span>{{end}}</h3>
                    {{if .CoverageData}}
                    <p>Отчет о покрытии кода тестами успешно сгенерирован.</p>
                    {{if .CoverageHTML}}
                    <div class="links">
                        <a href="{{.CoverageHTML}}" target="_blank">📊 Открыть HTML отчет о покрытии</a>
                    </div>
                    {{end}}
                    {{else}}
                    <p>Данные о покрытии кода тестами отсутствуют.</p>
                    {{end}}
                </div>
            </div>

            <!-- Gocyclo Section -->
            <div class="section">
                <h2>🔄 Цикломатическая сложность</h2>
                <div class="metric-card">
                    <h3>Статус: {{if .GocycloOutput}}<span class="status-ok">✓ Анализ завершен</span>{{else}}<span class="status-warning">⚠ Данные отсутствуют</span>{{end}}</h3>
                    {{if .GocycloOutput}}
                    <pre>{{.GocycloOutput}}</pre>
                    {{else}}
                    <p>Данные о цикломатической сложности отсутствуют.</p>
                    {{end}}
                </div>
            </div>
        </div>

        <footer>
            <p>Сгенерировано утилитой <strong>gokode</strong> | Все отчеты сохранены в директории <code>{{.MetricsDir}}</code></p>
        </footer>
    </div>
</body>
</html>
`
