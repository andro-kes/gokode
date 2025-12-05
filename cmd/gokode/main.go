package main

import (
	"fmt"
	"os"

	"github.com/andro-kes/gokode/worker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: gokode analyse <path>")
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "analyse" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintln(os.Stderr, "Usage: gokode analyse <path>")
		os.Exit(1)
	}

	path := "."
	if len(os.Args) > 2 {
		path = os.Args[2]
	}

	worker.Run(path)
}
