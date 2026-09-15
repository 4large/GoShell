package shell

import (
	"NoahShell/commands/cat"
	"NoahShell/commands/cd"
	"NoahShell/commands/echo"
	"NoahShell/commands/exit"
	"NoahShell/commands/grep"
	"NoahShell/commands/head"
	"NoahShell/commands/ls"
	"NoahShell/commands/mkdir"
	"NoahShell/commands/printhistory"
	"NoahShell/commands/pwd"
	"NoahShell/commands/rm"
	"NoahShell/commands/rmdir"
	"NoahShell/commands/touch"
	"NoahShell/internal/dispatch"
	"NoahShell/internal/history"
	"NoahShell/internal/pipeline"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

type ShellState struct {
	cwd          string
	history      []string
	lastExitCode int                  //Exit code of last cmd, so if last cmd resulted in error, its set to 1
	dispatcher   *dispatch.Dispatcher //command dispatch table reference
}

func Run() int {
	//register commands and build shell state (call dispatcher method)
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: could not get current directory")
		return 1
	}

	user, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to get user")
	}
	username := user.Username
	home := user.HomeDir

	History, err := history.Load(filepath.Join(home, ".gosh_history"))
	if err != nil && err != os.ErrNotExist {
		fmt.Fprintln(os.Stderr, "Error: failed to get history")
		return 1
	}

	sstate := ShellState{cwd, History, 0, dispatch.NewDispatcher()}

	//command registration
	sstate.dispatcher.RegisterCmd("cat", cat.Cat)
	sstate.dispatcher.RegisterCmd("ls", ls.Ls)
	sstate.dispatcher.RegisterCmd("cd", cd.Cd(&sstate.cwd, home))
	sstate.dispatcher.RegisterCmd("pwd", pwd.Pwd)
	sstate.dispatcher.RegisterCmd("echo", echo.Echo)
	sstate.dispatcher.RegisterCmd("touch", touch.Touch)
	sstate.dispatcher.RegisterCmd("exit", exit.Exit(&sstate.lastExitCode, &sstate.history, home))
	sstate.dispatcher.RegisterCmd("history", printhistory.MakeHistory(&sstate.history))
	sstate.dispatcher.RegisterCmd("mkdir", mkdir.Mkdir)
	sstate.dispatcher.RegisterCmd("rmdir", rmdir.Rmdir)
	sstate.dispatcher.RegisterCmd("rm", rm.Rm)
	sstate.dispatcher.RegisterCmd("grep", grep.Grep)
	sstate.dispatcher.RegisterCmd("head", head.Head)

	reader := bufio.NewReader(os.Stdin)

	//Create sig handle that properly handles ctrl c (save to history)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		for range sigs {
			history.Save(filepath.Join(home, ".gosh_history"), sstate.history)
			fmt.Fprintln(os.Stdout, "\nExiting shell, goodbye :)")
			os.Exit(0)
		}
	}()

	//Repl loop
	for {
		// print the prompt
		sstate.cwd = strings.Replace(sstate.cwd, home, "~", 1)
		fmt.Fprint(os.Stdout, username+":"+sstate.cwd+"$ ")

		// read a line of input
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				history.Save(filepath.Join(home, ".gosh_history"), sstate.history)
				fmt.Fprintln(os.Stdout, "\nExiting shell, goodbye :)")
				return sstate.lastExitCode
			} else {
				fmt.Fprintln(os.Stderr, "Error: failed to read user input")
				return 1
			}
		}

		// record it in history
		sstate.history = append(sstate.history, strings.TrimRight(input, "\n"))

		err2 := pipeline.Pipeline(input, os.Stdin, os.Stdout, sstate.dispatcher)
		if err2 != nil {
			sstate.lastExitCode = 1
			fmt.Fprintln(os.Stderr, err2)
		} else {
			sstate.lastExitCode = 0
		}
	}
}
