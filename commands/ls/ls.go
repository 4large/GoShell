package ls

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Ls(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("ls", flag.ContinueOnError)

	//Make flags for args
	flagA := fs.Bool("a", false, "Shows hidden files")
	flagL := fs.Bool("l", false, "Detailed files")
	flagH := fs.Bool("h", false, "Human readable sizes")
	flagN := fs.Bool("n", false, "enumerate UID/GID")
	flagR := fs.Bool("R", false, "Recursively probe subdirectories")

	//Instance variables
	fs.Parse(args[1:])
	flagBool := []bool{*flagA, *flagL, *flagH, *flagN, *flagR}
	isFlags := any(flagBool)
	files := fs.Args()

	writer := io.Writer(stdout)
	useColors := IsTerminal(stdout)

	//No args
	if len(files) == 0 {
		//get current directory
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		//POSIX equivalent opendir(3) readdir(3) where opendir returns open file and read dir reads the file.
		//There's another readdir call in ls but for sake of brevity I'm omitting it.
		entriesUnfiltered, err2 := os.ReadDir(dir)
		if err2 != nil {
			return err2
		}

		var entries []os.DirEntry
		if !flagBool[0] {
			entries = DirFilter(entriesUnfiltered)
		} else {
			entries = entriesUnfiltered
		}

		args := make([]string, len(entries))
		for i, entry := range entries {
			if entry == nil {
				continue
			}
			args[i] = filepath.Join(dir, entry.Name())
		}

		//Strips the leading . from the files for proper sorting
		sort.Slice(args, func(i, j int) bool {
			nameI := strings.TrimPrefix(filepath.Base(args[i]), ".")
			nameJ := strings.TrimPrefix(filepath.Base(args[j]), ".")
			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})

		//Append . and .. if -a
		if flagBool[0] {
			args = append([]string{dir + "/.", dir + "/.."}, args...)
		}

		if isFlags {
			fmt.Fprintln(writer, ".:")
			FlagsLs(writer, args, useColors, flagBool)
		} else {
			SimpleLS(writer, args, useColors)
		}

	} else {
		//sort files by files, directories before passing them in to be printed
		var fileArgs []string
		var directories []string
		for _, file := range files {
			//POSIX equivalent lstat(2), returns file metadata and doesn't follow symlinks, only metadata about link itself.
			info, err := os.Lstat(file)
			if err != nil {
				return err
			}

			if info.IsDir() {
				directories = append(directories, file)
			} else {
				fileArgs = append(fileArgs, file)
			}
		}

		sort.Slice(fileArgs, func(i, j int) bool {
			nameI := strings.TrimPrefix(filepath.Base(fileArgs[i]), ".")
			nameJ := strings.TrimPrefix(filepath.Base(fileArgs[j]), ".")
			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})

		sort.Slice(directories, func(i, j int) bool {
			nameI := strings.TrimPrefix(filepath.Base(directories[i]), ".")
			nameJ := strings.TrimPrefix(filepath.Base(directories[j]), ".")
			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})

		files = append(fileArgs, directories...)

		//files in the input NOTE! (.) means use current working directory
		for _, str := range files {
			fileInfo, err := os.Lstat(str)
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				//If we're in a pipe, we dont want to print dir name
				if _, ok := stdout.(*io.PipeWriter); !ok {
					fmt.Fprintln(writer, str+":")
				}
				dirRoot, err2 := os.ReadDir(str)
				var dirRootFilter []os.DirEntry
				//Filter only if flag a is not inputted
				if !flagBool[0] {
					dirRootFilter = DirFilter(dirRoot)
				} else {
					dirRootFilter = dirRoot
				}

				dirFiles := make([]string, len(dirRootFilter))
				if err2 != nil {
					return err2
				}

				for i, file := range dirRootFilter {
					dirFiles[i] = filepath.Join(str, file.Name())
				}

				sort.Slice(dirFiles, func(i, j int) bool {
					nameI := strings.TrimPrefix(filepath.Base(dirFiles[i]), ".")
					nameJ := strings.TrimPrefix(filepath.Base(dirFiles[j]), ".")
					return strings.ToLower(nameI) < strings.ToLower(nameJ)
				})

				dir, _ := os.Getwd()
				if flagBool[0] {
					dirFiles = append([]string{dir + "/.", dir + "/.."}, dirFiles...)
				}

				if isFlags {
					FlagsLs(writer, dirFiles, useColors, flagBool)
				} else {
					SimpleLS(writer, dirFiles, useColors)
				}

				fmt.Fprintln(writer, "")
			} else {
				var singleArr []string
				singleArr = append(singleArr, str)
				if isFlags {
					FlagsLs(writer, singleArr, useColors, flagBool)
				} else {
					SimpleLS(writer, singleArr, useColors)
				}
			}
		}
	}

	return nil
}

// Checks for flags in the input
func any(flags []bool) bool {
	for _, v := range flags {
		if v {
			return true
		}
	}
	return false
}
