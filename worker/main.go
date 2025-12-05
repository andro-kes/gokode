package worker

import (
	"fmt"
	"os"
	fp "path/filepath"

	"github.com/andro-kes/gokode/worker/analyser"
	"github.com/andro-kes/gokode/worker/jsoner"
	"github.com/andro-kes/gokode/worker/parser"
	"github.com/andro-kes/gokode/worker/scanner"
)

func Run(path string) {
	// Get the current directory before changing to the target path
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: Failed to get current directory")
		return
	}

	// Create metrics directory if it doesn't exist
	metricsDir := fp.Join(cwd, "metrics")
	err = os.MkdirAll(metricsDir, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: Failed to create metrics directory")
		return
	}

	jsonFile, err := os.Create(fp.Join(metricsDir, "report.json"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: Failed to create report.json")
		return
	}
	jsn := jsoner.NewJson(jsonFile)

	err = os.Chdir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: No such directory")
		return
	}

	fileChan := make(chan *os.File)

	scn := scanner.NewScanner()
	scn.FileChan = fileChan
	scn.Jsoner = jsn

	prs := parser.NewParser()
	prs.FileChan = fileChan
	prs.Jsoner = jsn

	go scn.Scan()
	prs.Parse()

	jsn.Write()

	vetFile, err := os.Create(fp.Join(metricsDir, "vet.txt"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: Failed to create vet.txt")
		return
	}
	analyser.Analyse(vetFile)

	defer func() {
		jsonFile.Close()
		vetFile.Close()
	}()
}
