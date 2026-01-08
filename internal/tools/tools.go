package tools

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	GolangciLintVersion = "v1.60.3"
	GocycloVersion      = "v0.6.0"
)

// IsInstalled checks if a tool is installed and available in PATH
func IsInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// InstallGolangciLint installs golangci-lint at the specified version
func InstallGolangciLint() error {
	fmt.Printf("Installing golangci-lint %s...\n", GolangciLintVersion)
	cmd := exec.Command("go", "install", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", GolangciLintVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InstallGocyclo installs gocyclo at the specified version
func InstallGocyclo() error {
	fmt.Printf("Installing gocyclo %s...\n", GocycloVersion)
	cmd := exec.Command("go", "install", fmt.Sprintf("github.com/fzipp/gocyclo/cmd/gocyclo@%s", GocycloVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InstallAll installs all required tools
func InstallAll() error {
	fmt.Println("Installing required tools...")

	var failed bool

	if err := InstallGolangciLint(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing golangci-lint: %v\n", err)
		failed = true
	} else {
		fmt.Println("✓ golangci-lint installed")
	}

	if err := InstallGocyclo(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing gocyclo: %v\n", err)
		failed = true
	} else {
		fmt.Println("✓ gocyclo installed")
	}

	if failed {
		return fmt.Errorf("failed to install some tools")
	}

	fmt.Println("✓ All tools installed successfully")
	return nil
}
