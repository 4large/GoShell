package ls

import (
	"fmt"
	"io"
	"path/filepath"
)

type Color string

const (
	C_BLUE  = Color("\033[34m")
	C_GREEN = Color("\033[32m")
	C_RESET = Color("\033[0m")
)

// Prints appropriate entries in color, or colorless if regular file
func (c Color) ColorPrint(w io.Writer, s string) {
	fmt.Fprintf(w, "%s%s%s\n", c, filepath.Base(s), C_RESET)
}
