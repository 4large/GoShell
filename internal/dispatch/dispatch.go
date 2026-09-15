package dispatch

import (
	"fmt"
	"io"
	"strings"
)

type CommandFunc func(args []string, stdin io.Reader, stdout io.Writer) error

type Dispatcher struct {
	commands map[string]CommandFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		commands: make(map[string]CommandFunc),
	}
}

// NOTE (d *Dispatcher) is like self in python
func (d *Dispatcher) RegisterCmd(name string, cmd CommandFunc) {
	d.commands[name] = cmd
}

func (d *Dispatcher) CmdLookup(name string) (CommandFunc, error) {
	name = strings.TrimRight(name, "\n")
	cmd, ok := d.commands[name]
	if !ok {
		return nil, fmt.Errorf("error: %s command not found", name)
	}
	return cmd, nil
}
