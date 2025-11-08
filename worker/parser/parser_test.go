package parser

import (
	"os"
	"testing"
	fp "path/filepath"

	"github.com/andro-kes/gokode/worker/jsoner"
)

func Init(t *testing.T) (*parser, *os.File) {
	t.Helper()

	file, err := os.Create("./test/report.json")
	if err != nil {
		t.Error("Failed to open test report.json")
	}
	jsn := jsoner.NewJson(file)
	fileChan := make(chan *os.File)

	prs := NewParser()
	prs.Jsoner = jsn
	prs.FileChan = fileChan

	t.Log("Parser is ready")

	return prs, file
}

func Walk(t *testing.T, prs *parser) {
	t.Helper()

	os.Chdir("./test")
	fp.Walk("./", func(filepath string, info os.FileInfo, err error) error {
		if fp.Ext(filepath) == ".go" {
			file, _ := os.Open(fp.Join("./", filepath))
			prs.Jsoner.AddFileName(filepath)
			prs.FileChan <- file
		}
        
        return nil
    })
	close(prs.FileChan)
}

func TestParse(t *testing.T) {
	prs, file := Init(t)
	go Walk(t, prs)
	prs.Parse()
	prs.Jsoner.Write()
	defer file.Close()
}