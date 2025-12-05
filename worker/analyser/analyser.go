package analyser

import (
	"fmt"
	"os"
	"os/exec"
)

func Analyse(vetFile *os.File) {
	Vet(vetFile)
}

func Vet(file *os.File) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Stdout = file
	cmd.Stderr = file
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Println("Success")
}
