# GoShell

A POSIX-style command shell written from scratch in Go — ~2,000 lines, 14 built-in commands, standard library only. Linux.

## Features

- **14 built-in commands**: `cd`, `pwd`, `ls` (`-a -l -h -n -R`), `cat`, `grep`, `head` (`-n`), `touch`, `mkdir`, `rmdir`, `rm` (`-r`, with confirmation prompt), `echo`, `history`, `exit`
- **Pipelines** — commands chained with concurrent reader/writer streams, modeled on `pipe(2)`/`dup(2)`
- **Redirection** — `&gt;` and `&gt;&gt;`
- **Persistent history** — saved to `~/.gosh_history`, flushed on exit or Ctrl-C
- **Graceful signal handling** — Ctrl-C and EOF shut down cleanly and save history
- **Terminal-aware output** — colored `ls` when stdout is a TTY, plain when piped
- **Exit codes** — `exit N` propagates `N` to the parent process

## How it works

All commands share one signature, `func(args, stdin, stdout) error`, and are registered by name in a dispatcher. Each command just reads `stdin` and writes `stdout` — it never knows whether it's talking to the terminal, a file, or another command. That's what makes piping work: the pipeline engine rewires those streams between stages.

Commands that need shell state (`cd`, `exit`, `history`) receive it through closures over pointers — `cd` gets a `*cwd` so directory changes show up in the prompt. This keeps the shell→command dependency one-directional and avoids circular imports.

Each pipeline stage runs in its own goroutine, connected to the next by an `io.Pipe`. Producers close their writers when done; readers close only after all stages join. Careful lifecycle ordering (plus draining upstream readers) is what prevents deadlocks and premature-close panics.

## Run it

Requires Go and Linux.

```bash
git clone https://github.com/4large/GoShell.git
cd GoShell
go run main.go        # or: go build -o gosh . && ./gosh
```
## Sample session

```bash
user:~/GoShell/testing$ cat out.txt | grep line | grep two
line two
user:~/GoShell/testing$ echo Hello Go >> out.txt
user:~/GoShell/testing$ cat out.txt
line one
line two
Hello Go
```

## Project layout

```
main.go            entry point
internal/
  shell/           REPL loop, state, command registration, signal handling
  dispatch/        command registry
  pipeline/        parsing, stage wiring, pipe lifecycle
  redirect/        > and >> redirection
  history/         ~/.gosh_history load/save
commands/          one package per built-in command
```
