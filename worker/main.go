package worker

import (
	"fmt"
	"os"

	"github.com/andro-kes/gokode/worker/analyser"
	"github.com/andro-kes/gokode/worker/jsoner"
	"github.com/andro-kes/gokode/worker/parser"
	"github.com/andro-kes/gokode/worker/scanner"
)

func Run(path string) {
	jsonFile, err := os.Create("../metrics/report.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: Failed to create report.json")
	}
	jsn := jsoner.NewJson(jsonFile)

	err = os.Chdir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: No such directory")
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

	vetFile, err := os.Create("../metrics/vet.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gokode: Failed to create vet.txt")
	}
	analyser.Analyse(vetFile)

	defer func(){
		jsonFile.Close()
		vetFile.Close()
	}()
}