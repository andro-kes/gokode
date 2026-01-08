package tools

import (
	"testing"
)

func TestIsInstalled(t *testing.T) {
	// Test with a tool that should definitely exist
	if !IsInstalled("go") {
		t.Error("Expected 'go' to be installed")
	}

	// Test with a tool that definitely doesn't exist
	if IsInstalled("this-tool-definitely-does-not-exist-12345") {
		t.Error("Expected non-existent tool to not be found")
	}
}
