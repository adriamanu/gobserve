[![Actions Status](https://github.com/adriamanu/goverwatch/actions/workflows/test.yml/badge.svg)](https://github.com/adriamanu/goverwatch/actions)
[![codecov](https://codecov.io/gh/adriamanu/goverwatch/master/graph/badge.svg)](https://codecov.io/gh/adriamanu/goverwatch)

# Golang file watcher

## What is Goverwatch
It is useful in any development environment to have a hot-reloading system.<br>
Govewatch is a file watcher that will lookup for modifications on a bunch of files and execute a command whenever a file is modified.<br>
Specify a **command** to execute and a list of files to watch on and you are ready to go.<br>

## Usage 
| Supported flags |                                    Description                                    | Mandatory |
| :-------------: | :-------------------------------------------------------------------------------: | :-------: |
|      files      | files separated by spaces you want to look on <br> * && ** patterns are supported |    yes    |
|        c        |                   command to execute wrapped with double quotes                   |    yes    |
|     ignore      | files separated by spaces you want to ignore <br> * && ** patterns are supported  |    no     |

Here is an example:
```bash
goverwatch -files "**/**/**/**/**.go **/**/**/*.json *.go" -c "go build main.go" -ignore "_samples/b/b.json *.go""
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

## Download
```bash
curl -L -s https://github.com/adriamanu/goverwatch/releases/download/latest/goverwatch  --output goverwatch && chmod +x goverwatch
```