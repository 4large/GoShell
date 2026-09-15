package ls

import (
	"os"
	"path/filepath"
)

// removes any hidden files from the dir listing, note not exported
func DirFilter(entries []os.DirEntry) []os.DirEntry {
	//Count number of files to exclude that way we can size the array properly
	counter := len(entries)
	for _, entry := range entries {
		if filepath.Base(entry.Name())[0] == 0x2E && len(entry.Name()) != 1 {
			counter -= 1
		}
	}

	if counter == len(entries) {
		return entries
	}

	args := make([]os.DirEntry, counter)
	i := 0
	for _, entry := range entries {
		if filepath.Base(entry.Name())[0] != 0x2E || len(entry.Name()) == 1 {
			args[i] = entry
			i++
		}
	}
	//size the array then create the list
	return args
}
