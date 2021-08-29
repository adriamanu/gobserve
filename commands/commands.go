package commands

import (
	"goverwatch/process"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type CommandToExecute struct {
	command string
	args    []string
}

// ParseCmd extract command and args from the provided string in a format used by 'Command' function from 'exec' package.
func ParseCmd(commandToParse string) CommandToExecute {
	tokenizedCommand := strings.Split(commandToParse, " ")
	amountOfToken := len(tokenizedCommand)

	cmd := tokenizedCommand[0]
	args := []string{}

	if amountOfToken > 1 {
		for i := 1; i < amountOfToken; i++ {
			args = append(args, tokenizedCommand[i])
		}
	}

	return CommandToExecute{cmd, args}
}

// ExecCmd execute the command and handle process.
func ExecCmd(c CommandToExecute) {
	process.KillRunningProcess()

	cmd := exec.Command(c.command, c.args...)
	// set the same pgid on childrens
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// redirect command errors and output to standard output and errors
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	process.RunningProcess = cmd.Process
}
