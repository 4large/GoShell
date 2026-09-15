package pwd

import (
	"fmt"
	"io"
	"os"
)

func Pwd(args []string, stdin io.Reader, stdout io.Writer) error {
	//POSIX equivalent getcwd(3), gets application cwd
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, dir)
	return nil
}
