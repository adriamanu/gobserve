package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var configFlag = flag.String("conf", "", "configuration file")
var commandFlag = flag.String("c", "", "command to execute")

type fileToWatch struct {
	filePath         string
	modificationDate int64
}

var filesToObserve []fileToWatch

func declareFilesToObserve(files []string) {
	for i := range files {
		filesToObserve = append(filesToObserve, fileToWatch{files[i], getFileModificationDate(files[i])})
	}
}

func getFileModificationDate(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	modifiedDate := fileInfo.ModTime().UnixNano()
	return modifiedDate
}

func fileHasBeenModified(originalModificationDate, modificationDate int64) bool {
	if originalModificationDate != modificationDate {
		return true
	}
	return false
}

func watchForModification(filesToWatch []fileToWatch, commandToExecute string) {
	for {
		for i := range filesToWatch {
			originalModificationDate := filesToWatch[i].modificationDate
			modificationDateToCheck := getFileModificationDate(filesToWatch[i].filePath)
			if fileHasBeenModified(originalModificationDate, modificationDateToCheck) {
				fmt.Printf("File %s has been modified\n", filesToWatch[i].filePath)
				filesToWatch[i].modificationDate = modificationDateToCheck
				execCmd(commandToExecute)
			}
		}
		time.Sleep(100 * time.Millisecond) // pause for 100 ms
	}
}

func parseCmd(commandToParse string) (string, []string) {
	tokenizedCommand := strings.Split(commandToParse, " ")
	amountOfToken := len(tokenizedCommand)
	cmd := tokenizedCommand[0]

	args := []string{}
	if amountOfToken > 1 {
		for i := 1; i < amountOfToken; i++ {
			args = append(args, tokenizedCommand[i])
		}
	}
	return cmd, args
}

func execCmd(commandToExecute string) {
	command, args := parseCmd(commandToExecute)
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	declareFilesToObserve([]string{"test.txt", "test2.txt"})
	watchForModification(filesToObserve, *commandFlag)
}
