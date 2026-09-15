package mkdir

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Creates directories of name args, similar to touch. Created with 0755 perms
// If directory already exits, print err (non fatal)
func Mkdir(args []string, stdin io.Reader, stdout io.Writer) error {
	dirs := args[1:]
	for _, dir := range dirs {
		cwd, err := ResolvePath(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		//POSIX equivalent mkdir(2), creates directory at specified path with specified octal perms.
		err2 := os.Mkdir(cwd, 0755)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, "Error: could not access directory or directory already exists: "+cwd)
		}
	}

	return nil
}

func ResolvePath(dir string) (string, error) {
	//Abs path
	if strings.Contains(dir, "/") {
		switch {
		case dir[0] == '/':
			return dir, nil
		case dir[:2] == "~/":
			home, err := os.UserHomeDir()
			if err != nil {
				return "", errors.New("error: failed to get home directory")
			}
			return filepath.Join(home, dir[2:]), nil
		default:
			cwd, err := os.Getwd()
			if err != nil {
				return "", errors.New("error: failed to get working directory")
			}
			return filepath.Join(cwd, dir), nil
		}
	}

	//Relative path (use cwd)
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.New("error: failed to get working directory")
	}
	return filepath.Join(cwd, dir), nil
}
