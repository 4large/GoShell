package exit

import (
	"NoahShell/internal/history"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// close stdout so the program exits the same way eof does
// Registers the function as a closure so we can edit and handle exit codes.
func Exit(exitCode *int, cmdhist *[]string, home string) func(args []string, stdin io.Reader, stdout io.Writer) error {
	return func(args []string, stdin io.Reader, stdout io.Writer) error {
		if len(args) == 1 {
			*exitCode = 0
		} else {
			code, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "exit: %s: invalid exit code\n", args[1])
				os.Exit(1)
			}
			*exitCode = code
		}
		history.Save(filepath.Join(home, ".gosh_history"), *cmdhist)
		fmt.Fprintln(os.Stdout, "Exiting shell, goodbye :)")
		//POSIX syscall exit(2), terminates process with the exit code, which watching processes may query.
		os.Exit(*exitCode)
		return nil
	}
}
