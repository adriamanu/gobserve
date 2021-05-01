package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Custom structures
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

// this slice if filled with globbed files
var filesToWatch []fileToWatch

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
	if runningProcess != nil {
		runningProcess.Signal(syscall.SIGKILL)
		// by using minus operator Kill will send signal to process group id
		// by doing so it will also kill process childrens
		err := syscall.Kill(-runningProcess.Pid, syscall.SIGKILL)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := exec.Command(c.command, c.args...)
	// Set the same pgid on childrens
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// The command will create a child process
	// If you use Kill() it will kill this process but not his childrens
	// His childrens will then be sons of INIT (PID 1) which can lead to unwanted scenarios
	// To prevent that we will kill all his children processes before killing him
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
	fmt.Println("Watching on", glob)
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
				fmt.Printf("File %s has been modified\n", filesToWatch[i].filePath)
				filesToWatch[i].modificationDate = modificationDateToCheck
				execCmd(cmd)
			}
		}
		// pause for 100 ms
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	flag.Parse()
	checkFlags()

	f := retrieveFilesToWatch(*filesFlag)
	declareFilesToWatch(f)

	cmd := parseCmd(*commandFlag)
	execCmd(cmd)
	watch(filesToWatch, cmd)
}
