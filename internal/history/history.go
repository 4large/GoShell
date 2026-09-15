package history

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Load history from text file
func Load(path string) ([]string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, errors.New("Error: failed to follow path")
	}

	//POSIX syscall open(2), returns open fd.
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Error: failed to open history file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	history := []string{}
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("Error: failed to read history file")
	}
	return history, nil
}

// write all entries to history file, if it does not exist create it
// file saved at ~/.gosh_history
func Save(path string, entries []string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.New("Error: failed to create history file")
	}
	defer f.Close()

	for _, entry := range entries {
		fmt.Fprintln(f, entry)
	}
	return nil
}
