package echo

import (
	"fmt"
	"io"
	"strings"
)

func Echo(args []string, stdin io.Reader, stdout io.Writer) error {
	outputs := args[1:]

	sb := strings.Builder{}
	for _, output := range outputs {
		sb.WriteString(output)
		sb.WriteString(" ")
	}
	fmt.Fprintln(stdout, sb.String())

	return nil
}
