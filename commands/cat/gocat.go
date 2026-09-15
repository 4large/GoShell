package cat

import (
	"flag"
	"io"
)

// Integrated from gocat
func Cat(args []string, stdin io.Reader, stdout io.Writer) error {
	//Read from args
	fs := flag.NewFlagSet("cat", flag.ContinueOnError)

	//Make flags for any method
	flagN := fs.Bool("n", false, "number output lines")
	flagB := fs.Bool("b", false, "number non empty output lines")
	flagS := fs.Bool("s", false, "Squeeze empty lines")
	flagE := fs.Bool("E", false, "display $ before each new line")
	flagT := fs.Bool("T", false, "Display tab character as ^I")
	flagV := fs.Bool("v", false, "Show non printing characters")
	flaglE := fs.Bool("e", false, "-vE shortcut")
	flaglT := fs.Bool("t", false, "-vT shortcut")
	flagA := fs.Bool("A", false, "-vET shortcut")

	fs.Parse(args[1:])

	flags := []bool{*flagN, *flagB, *flagS, *flagE, *flagT, *flagV, *flaglE,
		*flaglT, *flagA}

	files := fs.Args()

	//Decide which cat to use
	//NOTE! empty flag (-) is treated as false for the any function
	isFlags := any(flags)
	if isFlags {
		err := CatWFlags(files, flags, stdin, stdout)
		if err != nil {
			return err
		}
	} else {
		err := DumbCat(files, stdin, stdout)
		if err != nil {
			return err
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
