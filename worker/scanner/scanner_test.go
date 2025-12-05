package scanner

import (
	"os"
	"testing"

	"github.com/andro-kes/gokode/worker/jsoner"
)

func Init(t *testing.T) (*scanner, *os.File) {
	file, err := os.Create("./test/report.json")
	if err != nil {
		t.Errorf("gokode: %s", "Failed to create test reposrt.json")
	}

	jsn := jsoner.NewJson(file)
	fileChan := make(chan *os.File)

	scn := NewScanner()
	scn.Jsoner = jsn
	scn.FileChan = fileChan

	t.Log("Scanner is ready")

	return scn, file
}

func TestScan(t *testing.T) {
	scn, jsonFile := Init(t)

	os.Chdir("./test")
	go scn.Scan()

	for file := range scn.FileChan {
		t.Logf("%s: Done", file.Name())
		file.Close()
	}

	scn.Jsoner.Write()

	defer func() {
		jsonFile.Close()
	}()
}
