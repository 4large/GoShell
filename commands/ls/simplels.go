package ls

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// A simple ls for when no flags provided
func SimpleLS(w io.Writer, args []string, useColor bool) {
	for _, i := range args {
		info, err := os.Lstat(i)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		mode := info.Mode()

		//print with color, else print normally
		if useColor {
			if info.IsDir() {
				C_BLUE.ColorPrint(w, i)
			} else if mode.IsRegular() && (mode&0111) != 0 {
				C_GREEN.ColorPrint(w, i)
			} else {
				C_RESET.ColorPrint(w, i)
			}
		} else {
			//for piped inputs
			fmt.Fprintln(w, filepath.Base(i))
		}
	}
}
