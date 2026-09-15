package cat

import (
	"bufio"
	"io"
	"os"
)

// Flagged cat, setup writer and state, read each file and stdin if neccessary.
// read in through reader.go. preserve file state.
func CatWFlags(files []string, flags []bool, stdin io.Reader, stdout io.Writer) error {
	//If no files are given, read from stdin, apply flags
	writer := bufio.NewWriter(stdout)
	defer writer.Flush()

	lineNum := 1
	newLine := true
	var r io.Reader

	//if files is empty, get user input and then apply flags
	if len(files) == 0 {
		files = []string{"-"}
	}

	for _, i := range files {
		if i == "-" {
			r = stdin
		} else {
			//POSIX Equivalent open(2), returns open fd
			f, err := os.Open(i)
			if err != nil {
				return err
			}
			defer f.Close()
			r = f
		}

		err := reader(r, writer, flags, &lineNum, &newLine)
		if err != nil {
			return err
		}

		writer.Flush()
	}
	return nil
}
