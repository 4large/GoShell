package cat

import (
	"errors"
	"io"
	"os"
)

// flagless input, turns user input into output
// NOTE! (-) is passed in under "files"
func DumbCat(files []string, stdin io.Reader, stdout io.Writer) error {
	//handle no args, read from stdin
	if len(files) == 0 {
		_, err := io.Copy(stdout, stdin)

		//Sometimes, io.Copy throws an error when its upstream pipe is closed after the upstream pipe
		//Finishes sending its data so that error is irrelevant to us so we'll ignore it
		if errors.Is(err, io.ErrClosedPipe) {
			return nil
		}
		return err
	}

	for i := range files {
		if files[i] == "-" {
			//POSIX equivalent read(2), reads from a file descriptor.
			_, err := io.Copy(stdout, stdin)
			if err != nil {
				return err
			}
		} else {
			r, err := os.Open(files[i])
			if err != nil {
				return err
			}
			defer r.Close()

			//Check if dir before making copy
			info, err := r.Stat()
			if err != nil {
				r.Close()
				return err
			}
			if info.IsDir() {
				r.Close()
				return errors.New("Error: " + files[i] + " is a directory")
			}

			_, err2 := io.Copy(stdout, r)
			if err2 != nil {
				return err2
			}
		}

	}

	return nil
}
