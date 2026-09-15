package cd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
Cd uses a closure which has a reference to the shell state with the return function matching
command signature.

Since Go was being a royal pain, i could not import the shell or dispatch packages
for the shell state or CommandFunc signature so i have to pass in a pointer to
cwd and manually declare the types of the command func return signature.
*/
func Cd(cwd *string, home string) func(args []string, stdin io.Reader, stdout io.Writer) error {
	return func(args []string, stdin io.Reader, stdout io.Writer) error {
		//Replace ~ in cwd with home/user
		*cwd = strings.Replace(*cwd, "~", home, 1)

		//If the first char of our argument is / or ~ (or is empty) it means we're using an absolute path
		//not a relative one so we handle those base cases here
		var dir string
		switch {
		case len(args) == 1 || args[1] == "~":
			dir = home
		case args[1][0] == '~':
			dir = filepath.Join(home, args[1][1:])
		case args[1][0] == '/':
			dir = args[1]
		default:
			dir = filepath.Join(*cwd, args[1])
		}

		//POSIX equivalent chdir(2), changes cwd.
		err := os.Chdir(dir)
		if err != nil {
			return err
		}

		*cwd = dir
		return nil
	}
}
