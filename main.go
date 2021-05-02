package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// STRUCTURES
type fileToWatch struct {
	// location of the file
	filePath string
	// last time the file has been modified
	// if this value change the watcher will trigger a custom command
	modificationDate int64
}

type commandToExecute struct {
	command string
	args    []string
}

// VARS
var Reset = "\033[0m"
var Yellow = "\033[33m"

// this slice if filled with globbed files
var filesToWatch []fileToWatch

// currently running sub process
var runningProcess *os.Process

// FLAGS
var configFlag = flag.String("conf", "", "configuration file")
var ignoreFlag = flag.String("ignore", "", "regex of files to ignore")
var filesFlag = flag.String("files", "", "files to watch")
var commandFlag = flag.String("c", "", "command to execute")

func checkFlags() {
	errMsg := "Flag %s is mandatory"
	if *filesFlag == "" {
		log.Fatal(fmt.Errorf(errMsg, "-files"))
	}
	if *commandFlag == "" {
		log.Fatal(fmt.Errorf(errMsg, "-c"))
	}
}

// PROCESS MANAGEMENT
func killRunningProcess() {
	if runningProcess != nil {
		runningProcess.Signal(syscall.SIGKILL)
		// by using minus operator Kill will send signal to process group id
		// by doing so it will also kill process childrens
		err := syscall.Kill(-runningProcess.Pid, syscall.SIGKILL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func catchSigTerm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Printf(Yellow + "\n\nAnd now my watch is ended.\n" + Reset)
		killRunningProcess()
		os.Exit(0)
	}()
}

// COMMANDS
func parseCmd(commandToParse string) commandToExecute {
	tokenizedCommand := strings.Split(commandToParse, " ")
	amountOfToken := len(tokenizedCommand)
	cmd := tokenizedCommand[0]

	args := []string{}
	if amountOfToken > 1 {
		for i := 1; i < amountOfToken; i++ {
			args = append(args, tokenizedCommand[i])
		}
	}
	return commandToExecute{cmd, args}
}

func execCmd(c commandToExecute) {
	killRunningProcess()

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

	runningProcess = cmd.Process
}

// FILES
func retrieveFilesToWatch(pattern string) []string {
	glob, err := filepath.Glob(pattern)

	if err != nil {
		log.Fatal(err)
	}
	return glob
}

// we glob a list of files to watch and pass them to this function
// it add to our slice
func declareFilesToWatch(files []string) {
	for i := range files {
		filesToWatch = append(filesToWatch, fileToWatch{files[i], getFileModificationDate(files[i])})
	}
}

func getFileModificationDate(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.ModTime().UnixNano()
}

func fileHasBeenModified(originalModificationDate, modificationDate int64) bool {
	if originalModificationDate != modificationDate {
		return true
	}
	return false
}

func watch(filesToWatch []fileToWatch, cmd commandToExecute) {
	for {
		for i := range filesToWatch {
			originalModificationDate := filesToWatch[i].modificationDate
			modificationDateToCheck := getFileModificationDate(filesToWatch[i].filePath)
			if fileHasBeenModified(originalModificationDate, modificationDateToCheck) {
				fmt.Printf(Yellow+"File %s has been modified\n"+Reset, filesToWatch[i].filePath)
				filesToWatch[i].modificationDate = modificationDateToCheck
				execCmd(cmd)
			}
		}
	}
}

func main() {
	catchSigTerm()
	flag.Parse()
	checkFlags()

	fmt.Printf(Yellow + "And now my watch begins. It shall not end until my death.\n\n" + Reset)

	f := retrieveFilesToWatch(*filesFlag)
	declareFilesToWatch(f)
	fmt.Printf(Yellow+"watching on %s\n\n"+Reset, f)

	cmd := parseCmd(*commandFlag)

	execCmd(cmd)
	watch(filesToWatch, cmd)
}
