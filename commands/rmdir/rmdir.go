package rmdir

import (
	"NoahShell/commands/mkdir"
	"fmt"
	"io"
	"os"
)

// Removes non empty directories
func Rmdir(args []string, stdin io.Reader, stdout io.Writer) error {
	dirs := args[1:]
	for _, dir := range dirs {
		//We can use the path resolution from mkdir as we're resolving to a directory
		cwd, err := mkdir.ResolvePath(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		//Check if directory is empty, if so, remove it
		entries, err := os.ReadDir(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: failed to read directory "+cwd)
			continue
		}
		if len(entries) > 0 {
			fmt.Fprintln(os.Stderr, "Error: directory "+cwd+" is not empty")
			continue
		}

		//POSIX syscall unlink(2), deletes file at specified directory.
		err2 := os.Remove(cwd)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, "Error: failed to remove directory "+cwd)
		}
	}

	return nil
}
