# gokode

A standalone Go CLI utility for comprehensive code analysis and quality metrics. `gokode` can be installed and run against any Go project, providing formatting, vetting, linting with fixes, cyclomatic complexity analysis, and test coverage reports.

## Features

- **Format**: Automatically format code with `gofmt`
- **Vet**: Run `go vet` analysis and save results
- **Lint**: Run `golangci-lint` with JSON output and optional auto-fix
- **Test**: Run Go tests with coverage reports (profile and HTML)
- **Cyclomatic Complexity**: Analyze code complexity with `gocyclo`
- **Tool Bootstrap**: Automatically install required tools if missing
- **Metrics Reports**: All outputs saved to `metrics/` directory in target project

## Installation

Install `gokode` globally using Go:

```bash
go install github.com/andro-kes/gokode/cmd/gokode@latest
```

This will install the `gokode` binary to your `$GOPATH/bin` (or `$GOBIN`). Make sure this directory is in your `PATH`.

Alternatively, clone and build locally:

```bash
git clone https://github.com/andro-kes/gokode.git
cd gokode
go build -o gokode ./cmd/gokode
# Optionally move to PATH: sudo mv gokode /usr/local/bin/
```

## Usage

### Quick Start

Run full analysis on your current project:

```bash
gokode analyse .
```

Run analysis on a specific project:

```bash
gokode analyse /path/to/your/project
```

### Available Commands

```bash
gokode <command> [path]
```

**Commands:**

- `analyse` - Run full analysis (fmt, vet, lint with fixes, test, coverage, gocyclo)
- `fmt` - Format code with `gofmt`
- `vet` - Run `go vet` and write output to `metrics/vet.txt`
- `lint` - Run `golangci-lint` and write JSON to `metrics/report.json`
- `lint-fix` - Run `golangci-lint` with `--fix` flag
- `test` - Run tests with `go test ./...`
- `coverage` - Run tests with coverage (creates `metrics/coverage.out` and `coverage.html`)
- `gocyclo` - Run cyclomatic complexity analysis (writes to `metrics/gocyclo.txt`)
- `tools` - Install required tools (`golangci-lint`, `gocyclo`)

**Arguments:**

- `path` - Target directory (default: current directory `.`)

### Examples

```bash
# Run full analysis on current directory
gokode analyse .

# Format code in a specific project
gokode fmt ./myproject

# Run linting with auto-fix
gokode lint-fix /path/to/project

# Generate test coverage report
gokode coverage .

# Check cyclomatic complexity
gokode gocyclo .

# Install required tools
gokode tools
```

### Output Files

All analysis reports are written to a `metrics/` directory in the target project:

- `metrics/report.json` - golangci-lint results in JSON format
- `metrics/vet.txt` - go vet output
- `metrics/coverage.out` - test coverage profile
- `metrics/coverage.html` - test coverage HTML report
- `metrics/gocyclo.txt` - cyclomatic complexity analysis

The `metrics/` directory is created automatically if it doesn't exist.

## Tool Dependencies

`gokode` requires the following tools, which will be automatically installed if not found:

- **golangci-lint** (v1.60.3) - Comprehensive Go linter
- **gocyclo** (latest) - Cyclomatic complexity analyzer

To manually install all required tools:

```bash
gokode tools
```

## Configuration

### golangci-lint Configuration

The tool respects `.golangci.yml` configuration files in your project. If present, golangci-lint will use your custom configuration. The default configuration includes:

- Enabled linters: errcheck, gosimple, govet, staticcheck, unused, gofmt, goimports, misspell, gocritic, revive, gosec
- JSON output format for `metrics/report.json`
- 5-minute timeout for analysis

### Timeout

All operations have a default timeout of 5 minutes to prevent hanging on large projects.

## Development

### Running Tests

```bash
go test ./... -v
```

### Building Locally

```bash
go build -o gokode ./cmd/gokode
```

### Running on gokode Itself

```bash
# Analyze the gokode project
go run ./cmd/gokode analyse .

# Or with the built binary
./gokode analyse .
```

## CI/CD Integration

Use `gokode` in your CI/CD pipelines for automated code quality checks:

```yaml
# GitHub Actions example
- name: Install gokode
  run: go install github.com/andro-kes/gokode/cmd/gokode@latest

- name: Run analysis
  run: gokode analyse .

- name: Upload metrics
  uses: actions/upload-artifact@v4
  with:
    name: metrics
    path: metrics/
```

## Architecture

The CLI is built with Go's standard library and executes external tools for analysis:

- **Main CLI**: Command parsing and orchestration
- **Format**: Executes `gofmt -w -s`
- **Vet**: Executes `go vet ./...`
- **Lint**: Executes `golangci-lint run --out-format json`
- **Test/Coverage**: Executes `go test` with coverage flags
- **Gocyclo**: Executes `gocyclo` for complexity analysis

The legacy `worker/` package contains the original analyzer implementation and remains available for backward compatibility.

## License

This project is open source and available under the MIT License.
