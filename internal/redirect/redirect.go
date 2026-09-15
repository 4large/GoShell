package redirect

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// At the sink, truncate (>) or append (>>) to a file. redirection is only supported at the sink in a pipeline
// redirection is supported in non piped inputs.
func Redirection(cmd string) io.Writer {
	if strings.Contains(cmd, ">>") {
		return appendFile(cmd)
	} else {
		return truncate(cmd)
	}
}

func truncate(cmd string) io.Writer {
	args := strings.Split(cmd, ">")
	stdout := strings.TrimSpace(args[1])

	//POSIX syscall open(2), returns open file descriptors with perms and modifier flags
	f, err := os.OpenFile(stdout, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: was unable to truncate file "+stdout)
		return os.Stdout
	}
	//*os.File satisfies io.Writer requirement so we can return that
	return f
}

func appendFile(cmd string) io.Writer {
	args := strings.Split(cmd, ">>")
	stdout := strings.TrimSpace(args[1])

	f, err := os.OpenFile(stdout, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: was unable to append file "+stdout)
		return os.Stdout
	}
	return f
}
