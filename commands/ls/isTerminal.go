//Danny Radosevich

//Re-writing ls command in Go
//check if writing to terminal

package ls

import (
	"io"
	"os"
)

func IsTerminal(stdout io.Writer) bool {
	f, ok := stdout.(*os.File)
	if !ok {
		return false
	}

	//Checks file descriptor to see if were in the terminal or a pipe or a redirection
	fd := f.Fd()
	return fd == 0 || fd == 1 || fd == 2
}
