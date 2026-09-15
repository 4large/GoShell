package touch

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Creates a file of name args[1] and the specified POSIX metadata IE file perms 0644
func Touch(args []string, stdin io.Reader, stdout io.Writer) error {

	files := args[1:]
	for _, file := range files {
		cwd, err := resolvePath(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: failed to create file at path "+cwd)
			continue
		}

		//POSIX equivalent open(2), returns open fd, passes in flags to create fd metadata IE O_CREATE creates file if it
		//doesn't exist.
		f, err := os.OpenFile(cwd, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: failed to create file "+file)
			continue
		}
		defer f.Close()

		//POSIX equivalent utimensat(2), updates file timestamps, most precise utime syscall
		os.Chtimes(file, time.Now(), time.Now())
	}

	return nil
}

func resolvePath(file string) (string, error) {
	if strings.Contains(file, "/") {
		switch {
		case file[0] == '/':
			return file, nil
		case file[:2] == "~/":
			home, err := os.UserHomeDir()
			if err != nil {
				return "", errors.New("error: failed to get home directory")
			}
			return filepath.Join(home, file[2:]), nil
		default:
			cwd, err := os.Getwd()
			if err != nil {
				return "", errors.New("error: failed to get working directory")
			}
			return filepath.Join(cwd, file), nil
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.New("error: failed to get working directory")
	}
	return filepath.Join(cwd, file), nil
}
