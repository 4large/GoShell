[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/NyJBsIjs)
# [Noah Johnson]
## [Assignment Final]
## [Submission date: 04/28/26]
## Worked with/sources 
* https://pkg.go.dev/io
* https://stackoverflow.com/questions/29581165/trimright-everything-after-delimiter
* https://pkg.go.dev/os#OpenFile
* https://go.dev/blog/pipelines (particularly helpful)
## Project Quirks/ Things that don't work
* I had some trouble with package configs because I needed to import the shell state to other packages but the she already imported those packages creating a circular dependency. I would end up just using the fields from the shell state I needed passing a pointer but I wasn't aware that I should create a class for the shell state that didn't import anything.
* I had lots of trouble with concurrency and pipe readers/writer either with premature closings, deadlocks, improper error handling causing undefined behaviors, panics. Everything has to be perfectly in sync and it meant managing a lot of variables.
* Some concepts were hard to grasp, I wasn't initially sure what each stage was meant to do and the names for the functions I felt were misleading even if that's what they are in an actual shell (for example I thought dispatch would exit the command after pipeline sort of created a super command that linked other by a pipe, lot of learning about each of the stages was time consuming).
## Design Overview
* The code lifecycle runs within the shell where you can execute multiple commands and errors or completions don't terminate your shell. We begin by creating a shellstate which has process metadata like cwd, but most importantly contains our dispatcher. We use our dispatcher to register commands by name and their functions (Commmandfunc is our signature, we can register commands in a map by their name so we just query dispatch for the command when we need it). We also have a goroutine listening on our signals for anything that would end the program for graceful handling such as SIGINT or EOF. The rest of our program is contained in our repl loop where the overlying loop is print prompt, read input, record it in history, and then pass it to the pipeline and handle errors it throws. All commands are dispatched and executed from pipeline whether or not they are piped. We look on input for pipes or redirects and handle that accordingly. For redirections, its pretty simple in that we change stdout to the file of redirection and we pass that into the final input. For piped inputs, we create a pipe between our 2 commands. cmd1 writes to its pipewriter, which is read by cmd2 pipereader. We have to be very careful here and make sure to close our writers when were done and close our readers after all goroutines have wrapped up. We run the last command on main thread so that it blocks on its reader and doesn't prematurely exit the pipeline state. This directly mimics the POSIX syscalls pipe(2) and dup(2) where pipe(2) creates a reader and writer, and dup(2) does a fork to create a child process downstream that reads in from its parents stdout. Occasionally, we will need to access or modify the shell state. This was tricky because our shell imports pretty much everything so nothing can really import the shell because that creates circular dependency and won't compile. What we end up doing is passing pointers to the shell state that way all modifications are reported back to the shell IE cd needs to change the current working directory and that change needs to be reflected in the prompt. If our current working directory and the prompt our out of sync that represents a crucial error so what we do is we pass a pointer to cwd via a closure, as to not violate our command signature. when cd changes directory, cwd is updated and the change if reflected within shell state so that it can show the appropriate directory in the prompt.
## POSIX Syscall Map
| Command   | POSIX Syscall(s)            | Go Standard Library Call                                  |
|-----------|-----------------------------|-----------------------------------------------------------|
| `cd`      | `chdir(2)`                  | `os.Chdir()`                                              |
| `pwd`     | `getcwd(3)`                 | `os.Getwd()`                                              |
| `ls`      | `opendir(3)`, `readdir(3)`  | `os.ReadDir()`                                            |
| `cat`     | `open(2)`, `read(2)`        | `os.Open()`, `io.Copy()`                                  |
| `touch`   | `open(2)`, `utimensat(2)`   | `os.OpenFile(O_CREATE\|O_WRONLY)`, `os.Chtimes()`         |
| `mkdir`   | `mkdir(2)`                  | `os.Mkdir()`                                              |
| `rm`      | `unlink(2)`, `nftw()`       | `os.Remove()`, `os.RemoveAll()`                           |
| `rmdir`   | `rmdir(2)`                  | `os.Remove()`                                             |
| `grep`    | `read(2)`, `write(2)`       | `bufio.Scanner`, `io.Reader`, `io.Writer`                 |
| `echo`    | `write(2)`                  | `fmt.Fprintf()`                                           |
| `history` | `open(2)`, `read(2)`        | `os.Open()`, `bufio.Scanner`                              |
| `exit`    | `exit(2)`                   | `os.Exit()`                                               |
| `head`    | `read(2)`, `write(2)`       | `bufio.Scanner`, `io.Writer`                              |
| `>`       | `open(2)` `O_WRONLY\|O_TRUNC`   | `os.OpenFile(O_WRONLY\|O_CREATE\|O_TRUNC)`             |
| `>>`      | `open(2)` `O_WRONLY\|O_APPEND`  | `os.OpenFile(O_WRONLY\|O_CREATE\|O_APPEND)`            |
| `pipe (|)`| `pipe(2)`, `dup(2)`         | `io.Pipe()`, goroutines                                  |
## Known limitations
* The shell is only about 2000 LOC which is a far cry from a real shell so naturally it has limitations. Most prevelant, it assumes the user will use the commands in a fashion that makes sense. There is error handling, however, redirection only works on the final stage because there isn't a realistic scenario where you would need to use it in a pipe because that doesn't make a lot of sense. The shell will split on spaces rather than processing string literals for example, echo "Hello world" takes in the args 'echo', '"hello', and 'world"'. The shell also doesn't support quality of life features for example, you can't use up arrow key to use prev cmd from history. Also the commands are very simple in their implementations for example, mkdir doesn't support a recursively create directory flag.
## Build and run
* from project directory, use ./NoahShell or alternatively, you can use go run main.go. Only compatable with linux obviously.
## Sample Run

```bash
balls-johnson:~/GoProj/finalproject-4large/testing$ cat out.txt | grep line | grep two
line two 

balls-johnson:~/GoProj/finalproject-4large/testing$ echo Hello Go >> out.txt

balls-johnson:~/GoProj/finalproject-4large/testing$ cat out.txt
line one 
line two 
Hello Go 

balls-johnson:~/GoProj/finalproject-4large/testing$ cd ~
balls-johnson:~$ pwd
/home/balls-johnson

# (History was very long, so filtered with grep)
balls-johnson:~$ history | grep cd
4   cd testing
5   cd testing
6   cd ..
7   cd /
8   cd /
10  cd ..
12  cd ~
13  cd
15  cd commands
18  cd touch
29  cd testing
31  cd testing
33  cd testing
38  cd testing
40  cd ~
43  cd ~
46  cd testing
49  cd ~
52  cd testing
57  cd testing
73  cd testing
74  cd testing
76  cd dir1
88  cd testing
90  cd testing/dir1
93  cd ..
97  cd /
98  cd home
99  cd balls-johnson
102 cd /
107 cd testing
109 cd testing
113 cd testing
114 cd dir1
116 cd ndir
118 cd ../..
120 cd testing
125 cd ..
132 cd ~
134 cd GoProj/finalproject-4large
136 cd testing
139 cd testing
142 cd testing
145 cd testing
176 cd testing
210 cd testing
212 cd testing
215 cd testing
216 cd testing
221 cd testing | echo hello world | grep hello > foo.txt
223 cd ~
224 cd /tmp
225 cd
230 cd GoProj/finalproject-4large/testing
250 cd testing
260 cd testing
286 cd Testing
287 cd testing
294 cd ~
298 history | grep cd

balls-johnson:~/GoProj/finalproject-4large/testing$ rm -r file1.txt
Remove directory /home/balls-johnson/GoProj/finalproject-4large/testing/file1.txt and all contents? [y/N]: y

balls-johnson:~/GoProj/finalproject-4large/testing$ ls
file2.txt
out.txt
```


