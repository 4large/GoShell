package rm

import (
	"NoahShell/commands/mkdir"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// Remove a single file, if -r is specified, recursively remove directory and all it's contents
func Rm(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("ls", flag.ContinueOnError)

	flagR := fs.Bool("r", false, "Recursively delete directory")

	fs.Parse(args[1:])
	flagBool := *flagR
	files := fs.Args()

	for _, file := range files {
		cwd, err := mkdir.ResolvePath(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if flagBool {
			fmt.Fprint(os.Stdout, "Remove directory"+cwd+" and all contents? [y/N]:")
			reader := bufio.NewReader(stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: failed to read user input")
				continue
			}

			if input != "y\n" && input != "Y\n" {
				continue
			}

			//POSIX equivalent nftw(), recursively deletes all entries in a directory.
			err2 := os.RemoveAll(cwd)
			if err2 != nil {
				fmt.Fprintln(os.Stderr, "Error: failed to remove directory "+file)
			}
		} else {
			//no -r flag, check if its a directory and return err if it is
			fileinfo, err := os.Stat(cwd)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: failed to retrieve file data "+cwd)
				continue
			}

			if fileinfo.IsDir() {
				fmt.Fprintln(os.Stderr, "Error: the specified path "+cwd+" is a directory")
				continue
			}

			//POSIX syscall unlink(2), deletes file at specified directory.
			err2 := os.Remove(cwd)
			if err2 != nil {
				fmt.Fprintln(os.Stderr, "Error: failed to remove file "+file)
			}
		}
	}

	return nil
}
