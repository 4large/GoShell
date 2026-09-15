package ls

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Handles ls in the presence of flags
// Flags in ascending order: a, l, h, n, R
// NOTE! flags need to be applied first before the (potential) l or n calls
func FlagsLs(w io.Writer, files []string, useColor bool, flags []bool) {

	//check for n or l, else Normal ls with flags
	if flags[3] || flags[1] {
		printExpandedList(w, files, useColor, flags, ".")
	} else {
		printNormalList(w, files, useColor, flags, ".")
	}

}

// Called in the presence of the -n or -l flag
func printExpandedList(w io.Writer, files []string, useColor bool, flags []bool, header string) {
	//Traverse through all the entries to get the dimension of the entries metadata so we can format it nicely
	dimensions := [8]int{10, 0, 0, 0, 0, 3, 0, 0}
	for _, file := range files {
		//variables for extracting metadata
		info, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		stat := info.Sys().(*syscall.Stat_t)

		//Get Hard link count
		if len(fmt.Sprintf("%d", stat.Nlink)) > dimensions[1] {
			dimensions[1] = len(fmt.Sprintf("%d", stat.Nlink))
		}

		//Get owner name/UID group name/GID
		if flags[3] {
			if len(fmt.Sprintf("%d", os.Getuid())) > dimensions[2] {
				dimensions[2] = len(fmt.Sprintf("%d", os.Getuid()))
			}
			if len(fmt.Sprintf("%d", os.Getgid())) > dimensions[3] {
				dimensions[3] = len(fmt.Sprintf("%d", os.Getgid()))
			}
		} else {
			userName, _ := user.LookupId(fmt.Sprintf("%d", os.Getuid()))
			groupName, _ := user.LookupGroupId(fmt.Sprintf("%d", os.Getgid()))
			if len(userName.Username) > dimensions[2] {
				dimensions[2] = len(userName.Username)
			}
			if len(groupName.Name) > dimensions[3] {
				dimensions[3] = len(groupName.Name)
			}
		}

		//Get size, evaluate for -h
		if flags[2] {
			size := info.Size()
			switch {
			case size > 999:
				dimensions[4] = 4
			case size > 99 && dimensions[4] < 3:
				dimensions[4] = 3
			case size > 9 && dimensions[4] < 2:
				dimensions[4] = 2
			default:
				dimensions[4] = 1
			}
		} else {
			if len(fmt.Sprintf("%d", info.Size())) > dimensions[4] {
				dimensions[4] = len(fmt.Sprintf("%d", info.Size()))
			}
		}

		//Evaluate datetime as a unit
		modified := info.ModTime()
		cutoff := time.Now().AddDate(0, 0, -180)
		if modified.Before(cutoff) {
			//older than 180 days
			if dimensions[7] == 5 {
				continue
			} else {
				dimensions[7] = 4
			}
		} else {
			dimensions[7] = 5
		}

		if modified.Day() < 9 && dimensions[6] != 2 {
			dimensions[6] = 1
		} else {
			dimensions[6] = 2
		}
	}

	//Actually print out the files + file metadata
	for _, file := range files {
		//variables for extracting metadata
		info, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		stat := info.Sys().(*syscall.Stat_t)
		mode := info.Mode()

		//construct file permissions
		var sb strings.Builder
		if mode.IsDir() {
			sb.WriteString("d")
		} else if mode&os.ModeSymlink != 0 {
			sb.WriteString("l")
		} else {
			sb.WriteString("-")
		}

		// owner
		if mode&0400 != 0 {
			sb.WriteByte('r')
		} else {
			sb.WriteByte('-')
		}
		if mode&0200 != 0 {
			sb.WriteByte('w')
		} else {
			sb.WriteByte('-')
		}
		if mode&0100 != 0 {
			sb.WriteByte('x')
		} else {
			sb.WriteByte('-')
		}

		// group
		if mode&0040 != 0 {
			sb.WriteByte('r')
		} else {
			sb.WriteByte('-')
		}
		if mode&0020 != 0 {
			sb.WriteByte('w')
		} else {
			sb.WriteByte('-')
		}
		if mode&0010 != 0 {
			sb.WriteByte('x')
		} else {
			sb.WriteByte('-')
		}

		// others
		if mode&0004 != 0 {
			sb.WriteByte('r')
		} else {
			sb.WriteByte('-')
		}
		if mode&0002 != 0 {
			sb.WriteByte('w')
		} else {
			sb.WriteByte('-')
		}
		if mode&0001 != 0 {
			sb.WriteByte('x')
		} else {
			sb.WriteByte('-')
		}

		perms := sb.String()

		//get Nlink count
		Nlink := stat.Nlink

		//Get Username/uid and group name/gid
		uid := os.Getuid()
		gid := os.Getgid()

		var userName string
		var groupName string

		if flags[3] {
			// -n flag, print numeric uid/gid
			userName = fmt.Sprintf("%d", uid)
			groupName = fmt.Sprintf("%d", gid)
		} else {
			u, _ := user.LookupId(fmt.Sprintf("%d", uid))
			g, _ := user.LookupGroupId(fmt.Sprintf("%d", gid))
			userName = u.Username
			groupName = g.Name
		}

		//get size, apply -h flag
		size := float64(info.Size())
		counter := 0
		unitArr := [...]byte{0, 'K', 'M', 'G'}
		if flags[2] {
			for size > 1024 {
				size = size / 1024
				counter++
			}
		}

		//set it to a universal string, to handle flag and no flag
		var sizeStr string
		if unitArr[counter] == 0 {
			sizeStr = fmt.Sprintf("%d", int(size))
		} else if size >= 10 {
			sizeStr = fmt.Sprintf("%.0f%c", size, unitArr[counter]) // no decimal
		} else {
			sizeStr = fmt.Sprintf("%.1f%c", size, unitArr[counter]) // one decimal
		}

		//gets date last modified
		modified := info.ModTime().Local()
		cutoff := time.Now().AddDate(0, 0, -180)
		monthFull := modified.Month().String()
		month := monthFull[:3]
		day := modified.Day()

		var timeOrYear string
		if modified.Before(cutoff) {
			timeOrYear = fmt.Sprintf(" %d", modified.Year())
		} else {
			timeOrYear = modified.Format("15:04")
		}

		//print and format all the data
		if flags[2] {
			fmt.Fprintf(w, "%s %*d %-*s %-*s %4s %s %2d %s ",
				perms,
				dimensions[1], Nlink,
				dimensions[2], userName,
				dimensions[3], groupName,
				sizeStr,
				month,
				day,
				timeOrYear,
			)
		} else {
			fmt.Fprintf(w, "%s %*d %-*s %-*s %*d %s %2d %s ",
				perms,
				dimensions[1], Nlink,
				dimensions[2], userName,
				dimensions[3], groupName,
				dimensions[4], int(size),
				month,
				day,
				timeOrYear,
			)
		}

		//Print file normally
		if useColor {
			if info.IsDir() {
				C_BLUE.ColorPrint(w, file)
			} else if mode.IsRegular() && (mode&0111) != 0 {
				C_GREEN.ColorPrint(w, file)
			} else {
				C_RESET.ColorPrint(w, file)
			}
		} else {
			fmt.Fprintln(w, filepath.Base(file))
		}
	}

	//Handle recursion -R
	for _, file := range files {
		//do not recurse into . or .. dir
		if filepath.Base(file) == "." || filepath.Base(file) == ".." {
			continue
		}

		info, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		//Doesnt recurse into non directories or sym links
		if flags[4] && info.IsDir() {
			newHeader := header + "/" + filepath.Base(file)
			fmt.Fprintln(w, "")
			fmt.Fprintln(w, newHeader+":")

			entries, err := os.ReadDir(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				//Originally I had os.Exit here but that doesnt work in the context of a shell
				//And i cant return to the shell so i'll just continue instead
				continue
			}

			//Convert entries into a string array
			fileNames := make([]string, 0, len(entries))
			for _, entry := range entries {
				fileNames = append(fileNames, filepath.Join(file, entry.Name()))
			}

			if flags[0] {
				dir := filepath.Join(".", header[1:])
				fileNames = append([]string{dir + "/.", dir + "/.."}, fileNames...)
			}
			printExpandedList(w, fileNames, useColor, flags, newHeader)
		}
	}
}

// Normal flagged ls, handle -a (no explicit handling of a, handled elsewhere in the program) and -R
func printNormalList(w io.Writer, files []string, useColor bool, flags []bool, header string) {
	for _, file := range files {
		//Note, file is the file path as a string
		info, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		mode := info.Mode()

		//print with color, else print normally
		if useColor {
			if info.IsDir() {
				C_BLUE.ColorPrint(w, file)
			} else if mode.IsRegular() && (mode&0111) != 0 {
				C_GREEN.ColorPrint(w, file)
			} else {
				C_RESET.ColorPrint(w, file)
			}
		} else {
			//for piped inputs
			fmt.Fprint(w, filepath.Base(file))
		}
	}

	//Handle -R, call recursion after parent dir is printed
	for _, file := range files {
		//do not recurse into . or .. dir
		if filepath.Base(file) == "." || filepath.Base(file) == ".." {
			continue
		}

		info, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if flags[4] && info.IsDir() {
			newHeader := header + "/" + filepath.Base(file)
			fmt.Fprintln(w, "")
			fmt.Fprintln(w, newHeader+":")

			entries, err := os.ReadDir(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			//Convert entries into a string array
			fileNames := make([]string, 0, len(entries))
			for _, entry := range entries {
				fileNames = append(fileNames, filepath.Join(file, entry.Name()))
			}

			dir, _ := os.Getwd()
			if flags[0] {
				fileNames = append([]string{dir + "/.", dir + "/.."}, fileNames...)
			}
			printNormalList(w, fileNames, useColor, flags, newHeader)
		}
	}
}
