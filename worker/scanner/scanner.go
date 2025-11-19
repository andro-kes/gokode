package scanner

import (
	"fmt"
	"os"
	fp "path/filepath"

	"github.com/andro-kes/gokode/worker/jsoner"
)

type scanner struct {
	Jsoner *jsoner.Jsoner
	FileChan chan *os.File
}

func NewScanner() *scanner {
	return &scanner{}
}

func (s *scanner) Scan() {
	err := fp.Walk("./", func(filepath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

		if fp.Ext(filepath) == ".go" {
			file, err := os.Open(fp.Join("./", filepath))
			if err != nil {
				return err
			}

			s.Jsoner.AddFileName(filepath)

			select {
			case s.FileChan <- file:
			default:
				fmt.Fprintln(os.Stderr, "gokode: Channel was closed")
			}
		}
        
        return nil
    })

	defer close(s.FileChan)

	if err != nil {
		fmt.Fprintf(os.Stderr, "gokode: %s", err.Error())
		os.Exit(1)
	}
}