package grep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Grep takes input from stdin and matches strings on that. If its executed in isolation with no
// file it reads stdin.
// You can use it in like grep [string] [file], best implementation is read stdin and if
// stdin is empty then wait on input from stdin.
func Grep(args []string, stdin io.Reader, stdout io.Writer) error {
	args = args[1:]

	file := ""
	if len(args) < 1 {
		return errors.New("Error: malformed command 'grep'")
	} else if len(args) > 1 {
		file = args[1]
	}

	substr := args[0]
	f, err := os.Open(file)
	//POSIX equivalent, chained read(2) write(2) calls, scanner stores it into text field and then we write that out.
	var scanner *bufio.Scanner
	if err != nil {
		if len(args) > 1 {
			return errors.New("Error: directory " + args[1] + " doesn't exist")
		}
		scanner = bufio.NewScanner(stdin)
	} else {
		scanner = bufio.NewScanner(f)
	}
	defer f.Close()

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), substr) {
			fmt.Fprintln(stdout, scanner.Text())
		}
	}

	return nil
}
