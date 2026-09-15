package pipeline

import (
	"NoahShell/internal/dispatch"
	"NoahShell/internal/redirect"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// pipes commands together if given '|' as an input. commands in isolation will still be executed from the
// pipeline passed in via the repl loop but it will just take stdin from terminal and print stdout to terminal
// (or redirect)
func Pipeline(input string, stdin io.Reader, stdout io.Writer, dispatcher *dispatch.Dispatcher) error {
	//Dispatch all commands before we proceed
	cmds := parser(input)
	input = strings.Trim(input, "\n")
	var commandfuncs []dispatch.CommandFunc

	for _, cmd := range cmds {
		//get first word of string
		fields := strings.Fields(cmd)
		if len(fields) == 0 {
			return nil
		}
		cmd = fields[0]

		//Dispatch cmd
		cmdfn, err := dispatcher.CmdLookup(cmd)
		if err != nil {
			return err
		}

		commandfuncs = append(commandfuncs, cmdfn)
	}

	//Redirection only supported on last stage
	if strings.Contains(cmds[len(cmds)-1], ">") {
		f := redirect.Redirection(cmds[len(cmds)-1])
		if f == nil {
			stdout = os.Stdout
		} else {
			stdout = f
			file, _ := f.(*os.File)
			defer file.Close()
		}

		//trim anything after >, also do it for input since we pass input into commands in isolation
		tmp := strings.Split(cmds[len(cmds)-1], ">")
		cmds[len(cmds)-1] = strings.TrimSpace(tmp[0])
		input = strings.TrimSpace(strings.Split(input, ">")[0])
	}

	//pipe adjacent stages
	if len(commandfuncs) == 1 {
		return commandfuncs[0](strings.Fields(input), stdin, stdout)
	} else {
		return pipelineexc(commandfuncs, stdin, stdout, input)
	}
}

// Executed in the presence of "|"
func pipelineexc(cmds []dispatch.CommandFunc, stdin io.Reader, stdout io.Writer, input string) error {
	var wg sync.WaitGroup
	//stores our pipewriters so we can close it and there isnt confusion on which pw to close
	pipew := make([]*io.PipeWriter, len(cmds)-1)
	piper := make([]*io.PipeReader, len(cmds)-1)
	readers := make([]io.Reader, len(cmds))
	writers := make([]io.Writer, len(cmds))
	last := len(cmds) - 1

	//First stage needs stdin, last stage needs stdout, any intermediates need pr and pw
	readers[0] = stdin
	for i := 0; i < last; i++ {
		//Mimics pipe(2) syscall POSIX to create pipe readers and writers
		pr, pw := io.Pipe()
		pipew[i] = pw
		piper[i] = pr
		writers[i] = pw
		readers[i+1] = pr
	}
	writers[last] = stdout

	//build a 2d array of strings where each command is the split version on |. For example,
	//cat foo.txt | grep hello, each word is an element in the array and the arrays are indexed to be matched
	//With their command so when we run cat we call args[i] and gives the string array cat and foo.txt for cmd
	args := buildargs(input)

	for i := range cmds {
		wg.Add(1)
		go func(i int, cmd dispatch.CommandFunc) {
			defer wg.Done()
			if i < last {
				defer pipew[i].Close()
			}
			//Mimics dup(2) syscall POSIX to rewire stdin/stdout to pipereader/pipewriter
			err := cmd(args[i], readers[i], writers[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}(i, cmds[i])
	}

	wg.Wait()

	//Close the read pipes here to avoid premature closure
	for i := range piper {
		piper[i].Close()
	}

	//run last command on main thread so it blocks
	err := cmds[last](args[last], readers[last], writers[last])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	wg.Wait()

	return nil
}

// Parse on "|"
func parser(input string) []string {
	return strings.Split(input, "|")
}

func buildargs(input string) [][]string {
	args := [][]string{}
	parts := parser(input)
	for _, part := range parts {
		args = append(args, strings.Fields(part))
	}
	return args
}
