# gokode

Самостоятельная утилита CLI на Go для комплексного анализа кода и метрик качества. `gokode` может быть установлена и запущена для любого Go проекта, обеспечивая форматирование, проверку, линтинг с автоисправлениями, анализ цикломатической сложности и отчеты о покрытии тестами.

---

A standalone Go CLI utility for comprehensive code analysis and quality metrics. `gokode` can be installed and run against any Go project, providing formatting, vetting, linting with fixes, cyclomatic complexity analysis, and test coverage reports.

## Возможности / Features

- **Форматирование / Format**: Автоматическое форматирование кода с помощью `gofmt`
- **Проверка / Vet**: Запуск анализа `go vet` и сохранение результатов
- **Линтинг / Lint**: Запуск `golangci-lint` с форматированным JSON выводом и опциональным автоисправлением
- **Тестирование / Test**: Запуск Go тестов с отчетами о покрытии (профиль и HTML)
- **Цикломатическая сложность / Cyclomatic Complexity**: Анализ сложности кода с помощью `gocyclo`
- **HTML отчеты / HTML Reports**: Агрегированные HTML отчеты со всеми метриками (**НОВОЕ!**)
- **Автоустановка инструментов / Tool Bootstrap**: Автоматическая установка необходимых инструментов при отсутствии
- **Отчеты о метриках / Metrics Reports**: Все выходные данные сохраняются в директории `metrics/` целевого проекта

## Установка / Installation

Установить `gokode` глобально используя Go:

```bash
go install github.com/andro-kes/gokode/cmd/gokode@latest
```

Это установит бинарный файл `gokode` в ваш `$GOPATH/bin` (или `$GOBIN`). Убедитесь, что эта директория находится в вашем `PATH`.

Альтернативно, клонируйте и соберите локально:

```bash
git clone https://github.com/andro-kes/gokode.git
cd gokode
go build -o gokode ./cmd/gokode
# Опционально переместите в PATH: sudo mv gokode /usr/local/bin/
```

## Использование / Usage

### Быстрый старт / Quick Start

Запустить полный анализ текущего проекта:

```bash
gokode analyse .
```

Запустить анализ определенного проекта:

```bash
gokode analyse /path/to/your/project
```

### Доступные команды / Available Commands

```bash
gokode <command> [path]
```

**Команды / Commands:**

- `analyse` - Запустить полный анализ (fmt, vet, lint с автоисправлениями, test, coverage, gocyclo) и сгенерировать HTML отчет / Run full analysis (fmt, vet, lint with fixes, test, coverage, gocyclo) and generate HTML report
- `fmt` - Форматировать код с помощью `gofmt` / Format code with `gofmt`
- `vet` - Запустить `go vet` и записать вывод в `metrics/vet.txt` / Run `go vet` and write output to `metrics/vet.txt`
- `lint` - Запустить `golangci-lint` и записать форматированный JSON в `metrics/report.json` / Run `golangci-lint` and write pretty-printed JSON to `metrics/report.json`
- `lint-fix` - Запустить `golangci-lint` с флагом `--fix` / Run `golangci-lint` with `--fix` flag
- `test` - Запустить тесты с `go test ./...` / Run tests with `go test ./...`
- `coverage` - Запустить тесты с покрытием (создает `metrics/coverage.out` и `coverage.html`) / Run tests with coverage (creates `metrics/coverage.out` and `coverage.html`)
- `gocyclo` - Запустить анализ цикломатической сложности (записывает в `metrics/gocyclo.txt`) / Run cyclomatic complexity analysis (writes to `metrics/gocyclo.txt`)
- `tools` - Установить необходимые инструменты (`golangci-lint`, `gocyclo`) / Install required tools (`golangci-lint`, `gocyclo`)

**Аргументы / Arguments:**

- `path` - Целевая директория (по умолчанию: текущая директория `.`) / Target directory (default: current directory `.`)

### Примеры / Examples

```bash
# Запустить полный анализ текущей директории
# Run full analysis on current directory
gokode analyse .

# Форматировать код в определенном проекте
# Format code in a specific project
gokode fmt ./myproject

# Запустить линтинг с автоисправлением
# Run linting with auto-fix
gokode lint-fix /path/to/project

# Сгенерировать отчет о покрытии тестами
# Generate test coverage report
gokode coverage .

# Проверить цикломатическую сложность
# Check cyclomatic complexity
gokode gocyclo .

# Установить необходимые инструменты
# Install required tools
gokode tools
```

### Выходные файлы / Output Files

Все отчеты анализа записываются в директорию `metrics/` в целевом проекте:

All analysis reports are written to a `metrics/` directory in the target project:

- `metrics/report.json` - результаты golangci-lint в форматированном JSON / golangci-lint results in pretty-printed JSON format
- `metrics/report.html` - агрегированный HTML отчет со всеми метриками (**НОВОЕ!**) / aggregated HTML report with all metrics (**NEW!**)
- `metrics/vet.txt` - вывод go vet / go vet output
- `metrics/coverage.out` - профиль покрытия тестами / test coverage profile
- `metrics/coverage.html` - HTML отчет о покрытии тестами / test coverage HTML report
- `metrics/gocyclo.txt` - анализ цикломатической сложности / cyclomatic complexity analysis

Директория `metrics/` создается автоматически, если она не существует.

The `metrics/` directory is created automatically if it doesn't exist.

### HTML отчет / HTML Report

При запуске команды `analyse`, автоматически генерируется удобный для разработчика HTML отчет, который:

When running the `analyse` command, a developer-friendly HTML report is automatically generated that:

- Агрегирует все метрики в одном месте / Aggregates all metrics in one place
- Отображает результаты go vet, golangci-lint, покрытие тестами и цикломатическую сложность / Displays go vet, golangci-lint results, test coverage, and cyclomatic complexity
- Включает ссылки на сгенерированные артефакты (JSON отчеты, HTML покрытие) / Includes links to generated artifacts (JSON reports, HTML coverage)
- Использует современный, адаптивный дизайн с цветовым кодированием / Uses modern, responsive design with color coding
- Локализован на русском языке / Localized in Russian

Откройте `metrics/report.html` в браузере после запуска анализа для просмотра агрегированных результатов.

Open `metrics/report.html` in a browser after running analysis to view aggregated results.

## Зависимости инструментов / Tool Dependencies

`gokode` требует следующие инструменты, которые будут автоматически установлены при отсутствии:

`gokode` requires the following tools, which will be automatically installed if not found:

- **golangci-lint** (v1.60.3) - Комплексный линтер Go / Comprehensive Go linter
- **gocyclo** (v0.6.0) - Анализатор цикломатической сложности / Cyclomatic complexity analyzer

Для ручной установки всех необходимых инструментов:

To manually install all required tools:

```bash
gokode tools
```

## Конфигурация / Configuration

### Конфигурация golangci-lint / golangci-lint Configuration

Инструмент учитывает файлы конфигурации `.golangci.yml` в вашем проекте. Если они присутствуют, golangci-lint будет использовать вашу пользовательскую конфигурацию. Конфигурация по умолчанию включает:

The tool respects `.golangci.yml` configuration files in your project. If present, golangci-lint will use your custom configuration. The default configuration includes:

- Включенные линтеры / Enabled linters: errcheck, gosimple, govet, staticcheck, unused, gofmt, goimports, misspell, gocritic, revive, gosec
- Формат вывода JSON для `metrics/report.json` (с форматированием) / JSON output format for `metrics/report.json` (pretty-printed)
- Таймаут 5 минут для анализа / 5-minute timeout for analysis

### Таймаут / Timeout

Все операции имеют таймаут по умолчанию в 5 минут для предотвращения зависания на больших проектах.

All operations have a default timeout of 5 minutes to prevent hanging on large projects.

## Архитектура / Architecture

CLI построен с использованием стандартной библиотеки Go и выполняет внешние инструменты для анализа:

The CLI is built with Go's standard library and executes external tools for analysis:

- **Главный CLI / Main CLI** (`cmd/gokode`): Разбор команд и оркестрация / Command parsing and orchestration
- **Пакет Runner / Runner Package** (`internal/runner`): Логика выполнения команд для форматирования, проверки, линтинга, тестирования и анализа сложности / Command execution logic for format, vet, lint, test, and complexity analysis
- **Пакет Report / Report Package** (`internal/report`): Генерация HTML отчетов из файлов метрик / HTML report generation from metrics files
- **Пакет Tools / Tools Package** (`internal/tools`): Установка и управление инструментами / Tool installation and management

Каждая команда выполняет соответствующие внешние инструменты:

Each command executes the corresponding external tools:

- **Format**: Выполняет `gofmt -w -s` / Executes `gofmt -w -s`
- **Vet**: Выполняет `go vet ./...` / Executes `go vet ./...`
- **Lint**: Выполняет `golangci-lint run --out-format json` / Executes `golangci-lint run --out-format json`
- **Test/Coverage**: Выполняет `go test` с флагами покрытия / Executes `go test` with coverage flags
- **Gocyclo**: Выполняет `gocyclo` для анализа сложности / Executes `gocyclo` for complexity analysis

Устаревший пакет `worker/` содержит оригинальную реализацию анализатора и остается доступным для обратной совместимости.

The legacy `worker/` package contains the original analyzer implementation and remains available for backward compatibility.

## Разработка / Development

### Запуск тестов / Running Tests

```bash
go test ./... -v
```

### Локальная сборка / Building Locally

```bash
go build -o gokode ./cmd/gokode
```

### Запуск на самом gokode / Running on gokode Itself

```bash
# Анализ проекта gokode
# Analyze the gokode project
go run ./cmd/gokode analyse .

# Или с собранным бинарным файлом
# Or with the built binary
./gokode analyse .
```

## Интеграция CI/CD / CI/CD Integration

Используйте `gokode` в ваших CI/CD пайплайнах для автоматической проверки качества кода:

Use `gokode` in your CI/CD pipelines for automated code quality checks:

```yaml
# Пример для GitHub Actions
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

## Лицензия / License

Этот проект является открытым и доступен под лицензией MIT.

This project is open source and available under the MIT License.
