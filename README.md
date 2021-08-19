# Goverwatch

![goverwatch tests](https://github.com/adriamanu/goverwatch/actions/workflows/test.yml/badge.svg)
## Golang file watcher

Specify a command to execute and a list of files to watch on.<br>
If one of the file is modified, your command will be re-executed.<br>
It is useful in a development environment to have hot reloading.

```bash
echo "trigger watcher" > test.txt
go run main.go -c "go build main.go" -files test.txt
```

```bash
go run goverwatch.go -c "go run test/server.go" -files "*/*.go"
```

## Processes
`cmd.Start()` will create a child process.<br>
If you use the `Kill()` function it will kill this process but not his childrens.<br>
His childrens will then be sons of INIT (PID 1) which can lead to unwanted scenarios.<br>
To prevent that we are using process group in order to kill process and all of his childrens.<br>

To kill this process and his childrens we have to add them in the same process group.<br>
```go
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
```
We can then kill every process in the group by passing a negative integer to the `Kill()` function.<br>
```go
syscall.Kill(-runningProcess.Pid, syscall.SIGKILL)
```