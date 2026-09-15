package printhistory

import (
	"fmt"
	"io"
)

// Prints out history from ~/.gosh_history
func MakeHistory(cmdHist *[]string) func(args []string, stdin io.Reader, stdout io.Writer) error {
	return func(args []string, stdin io.Reader, stdout io.Writer) error {
		for i, entry := range *cmdHist {
			fmt.Fprintf(stdout, "%d %s\n", i+1, entry)
		}
		return nil
	}
}
