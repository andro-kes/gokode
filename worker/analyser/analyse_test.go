package analyser

import (
	"os"
	"testing"
)

func Init(t *testing.T) *os.File{
	t.Helper()

	file, err := os.Create("./test/vet.txt")
	if err != nil {
		t.Error("Failed to open test vet.txt")
	}

	t.Log("Vet file is ready")

	return file
}

func TestAnalyse(t *testing.T) {
	file := Init(t)
	os.Chdir("./test")
	Analyse(file)
}