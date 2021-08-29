package process

import (
	"fmt"
	"goverwatch/colors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// currently running sub process
var RunningProcess *os.Process

// KillRunningProcess send a Kill signal to children process
func KillRunningProcess() {
	if RunningProcess != nil {
		RunningProcess.Signal(syscall.SIGKILL)
		// by using minus operator Kill will send signal to process group id
		// by doing so it will also kill process childrens
		err := syscall.Kill(-RunningProcess.Pid, syscall.SIGKILL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CatchSignalsAndExit catch kill and terminal closing signals and clean exit the program
func CatchSignalsAndExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Printf(colors.Yellow + "\n\nAnd now my watch is over.\n" + colors.Reset)
		KillRunningProcess()
		os.Exit(0)
	}()
}
