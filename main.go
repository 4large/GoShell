package main

import (
	"NoahShell/internal/shell"
	"os"
)

func main() {
	exitcode := shell.Run()
	os.Exit(exitcode)
}
