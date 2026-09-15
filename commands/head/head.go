package head

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Prints first n lines of either a file or stdin
func Head(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("head", flag.ContinueOnError)

	flagN := fs.Bool("n", false, "Print first N lines of stdin, if no N print deafult is 10")

	fs.Parse(args[1:])
	flagBool := *flagN
	files := fs.Args()

	n, err := setn(flagBool, files)
	if err != nil {
		return err
	}

	if flagBool {
		return nflag(files, n, stdin, stdout)
	} else {
		return noflag(files, n, stdin, stdout)
	}
}

// flag
func nflag(files []string, n int, stdin io.Reader, stdout io.Writer) error {
	var scanner *bufio.Scanner
	var f *os.File
	var err error

	/*there is very different behavior for both head -n 5 and head -n 5 foo.txt. Our loop
	needs to process whether we have a file listed or not so it needs to work on 5 and then loop
	on stdin and if there are files, we need to continue to the next iteration on first pass.
	*/
	for i, file := range files {
		if len(files) < 2 {
			scanner = bufio.NewScanner(stdin)
		} else {
			//First pass with a file to output will have the file set to n IE head n 5 foo.txt,
			//file will be 5 so we want to continue onto foo
			if i == 0 {
				continue
			}
			f, err = os.Open(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.New("Error: failed to open file "+file))
				continue
			}
			defer f.Close()

			scanner = bufio.NewScanner(f)
		}

		stop := 0
		for scanner.Scan() {
			if stop >= n {
				break
			}
			fmt.Fprintln(stdout, scanner.Text())
			stop++
		}

	}

	//rest of upstream needs to be drained to prevent infinite blocking
	io.Copy(io.Discard, stdin)

	return nil
}

// flagless, this works pretty much the same as flagged except we check len < 1 instead of 2 since
// -n is not in the args anymore, also we dont need to skip the first arg when there are multiple.
func noflag(files []string, n int, stdin io.Reader, stdout io.Writer) error {
	var scanner *bufio.Scanner
	var f *os.File
	var err error

	if len(files) < 1 {
		scanner = bufio.NewScanner(stdin)

		stop := 0
		for scanner.Scan() {
			if stop >= n {
				break
			}
			fmt.Fprintln(stdout, scanner.Text())
			stop++
		}
	} else {
		for _, file := range files {
			f, err = os.Open(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.New("Error: failed to open file "+file))
				continue
			}
			defer f.Close()

			scanner = bufio.NewScanner(f)

			stop := 0
			for scanner.Scan() {
				if stop >= n {
					break
				}
				fmt.Fprintln(stdout, scanner.Text())
				stop++
			}
		}
	}

	//rest of upstream needs to be drained to prevent infinite blocking
	io.Copy(io.Discard, stdin)

	return nil
}

func setn(fl bool, files []string) (int, error) {
	//flagless, else flagged
	if !fl {
		return 10, nil
	} else {
		//If flagged we need at least one arg otherwise command is malformed
		if len(files) < 1 {
			return -1, errors.New("Error: malformed command, please specify length with -n flag")
		} else {
			n, err := strconv.Atoi(files[0])
			if err != nil {
				return -1, errors.New("Error: invalid length for n flag " + files[0])
			}
			return n, nil
		}
	}
}
