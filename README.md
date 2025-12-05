# gokode

A Go code analysis tool that provides metrics on code quality, including line counts, function counts, and linting results.

## Features

- Scan and analyze Go source files
- Generate metrics reports in JSON format
- Run `go vet` analysis
- Format, lint, and fix code issues
- Generate test coverage reports

## Installation

Clone the repository:
```bash
git clone https://github.com/andro-kes/gokode.git
cd gokode
```

Install dependencies (including golangci-lint v1.55.2):
```bash
make deps
```

Build the binary:
```bash
make build
```

Or build with Go directly:
```bash
go build -o gokode ./cmd/gokode
```

## Usage

### CLI Tool

Run analysis on the current directory:
```bash
./gokode analyse .
```

Run analysis on a specific directory:
```bash
./gokode analyse /path/to/your/project
```

Run without building (using go run):
```bash
go run ./cmd/gokode analyse .
```

The tool will generate two files in the `metrics/` directory:
- `metrics/report.json` - JSON report with code metrics
- `metrics/vet.txt` - Output from `go vet`

### Makefile Targets

The project includes a comprehensive Makefile with the following targets:

- `make help` - Show all available targets
- `make build` - Build the gokode binary
- `make run` - Build and run the tool on the current directory
- `make test` - Run all tests
- `make vet` - Run go vet
- `make fmt` - Format code with gofmt
- `make lint` - Run golangci-lint
- `make lint-fix` - Run golangci-lint with --fix
- `make deps` - Install dependencies including golangci-lint
- `make test-coverage` - Run tests with coverage and generate HTML report
- `make analyse` - Run full analysis: format, vet, and lint with fixes
- `make clean` - Clean build artifacts and metrics
- `make docker-build` - Build Docker image
- `make docker-run` - Build and run Docker container

### Full Analysis

To run a complete code analysis (format, vet, lint):
```bash
make analyse
```

This will:
1. Format your code with `gofmt`
2. Run `go vet` and write results to `metrics/vet.txt`
3. Run `golangci-lint` with auto-fix enabled
4. Generate a JSON report in `metrics/report.json`

## Configuration

### Linting Configuration

The project uses `.golangci.yml` for golangci-lint configuration. The configuration includes:
- Common linters: errcheck, gosimple, govet, staticcheck, unused, gofmt, goimports, misspell, gocritic, revive, gosec
- Output format: JSON
- Go version: 1.24

You can customize the linting rules by editing `.golangci.yml`.

## Development

Run tests:
```bash
make test
```

Generate test coverage:
```bash
make test-coverage
```

Format code:
```bash
make fmt
```

Lint code:
```bash
make lint
```

## Architecture

The project consists of several components:

- **Scanner**: Scans all .go files and sends them to the Parser through a channel
- **Parser**: Reads from the channel and parses files line-by-line to extract metrics (active lines of code, number of functions)
- **Analyser**: Analyzes the directory and collects metrics using commands like `go vet`
- **Jsoner**: Handles JSON output and report generation

For more details, see [worker/README.md](worker/README.md).
