package parser

import (
	"bufio"
	"os"
	"sync"

	"github.com/andro-kes/gokode/worker/jsoner"
)

const WORKERS int = 5

type parser struct {
	FileChan chan *os.File
	Jsoner *jsoner.Jsoner
	wg sync.WaitGroup
}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) Parse() {
	for range WORKERS {
		p.wg.Add(1)
		go p.parseWorker()
	}
	p.wg.Wait()
}

func (p *parser) parseWorker() {
	defer p.wg.Done()
	for file := range p.FileChan {
		scn := bufio.NewScanner(file)
		rows := 0
		for scn.Scan() {
			rows++
		}
		p.Jsoner.AddFileMetric(file.Name(), "number_of_rows", rows)
		file.Close()
	}
}