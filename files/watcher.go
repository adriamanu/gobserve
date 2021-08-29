package files

import (
	"fmt"
	"goverwatch/colors"
	"goverwatch/commands"
	"log"
	"os"
)

type fileToWatch struct {
	// filePath is the location of the file.
	filePath string
	// modificationDate is the last time this file has been modified.
	// if this value change the watcher will trigger a custom command.
	modificationDate int64
}

// filesToWatch is a slice filled with globbed files
var filesToWatch []fileToWatch

// DeclareFilesToWatch takes a list of files and add them to the global variable 'filesToWatch'
// with the filepath and the last modification date.
func DeclareFilesToWatch(files []string) {
	for i := range files {
		filesToWatch = append(filesToWatch, fileToWatch{files[i], getFileModificationDate(files[i])})
	}
}

// getFileModificationDate retrieve last modification date of a file thanks to 'Stat' function
// from 'os' package and returns an unix timestamp.
func getFileModificationDate(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.ModTime().UnixNano()
}

// fileHasBeenModified compare the modified date with the original one.
// true is returned if the file has been modified otherwise it returns false.
func fileHasBeenModified(originalModificationDate, modificationDate int64) bool {
	if originalModificationDate != modificationDate {
		return true
	}
	return false
}

// executeCommandOnFileModification will execute the provided command whenever a file is modified.
func executeCommandOnFileModification(cmd commands.CommandToExecute, files []fileToWatch) {
	for i := range filesToWatch {
		originalModificationDate := filesToWatch[i].modificationDate
		modificationDateToCheck := getFileModificationDate(filesToWatch[i].filePath)
		if fileHasBeenModified(originalModificationDate, modificationDateToCheck) {
			fmt.Printf(colors.Yellow+"File %s has been modified\n"+colors.Reset, filesToWatch[i].filePath)
			filesToWatch[i].modificationDate = modificationDateToCheck
			commands.ExecCmd(cmd)
		}
	}
}

// Watch is an infinite loop that executes 'executeCommandOnFileModification' function.
func Watch(cmd commands.CommandToExecute) {
	for {
		executeCommandOnFileModification(cmd, filesToWatch)
	}
}
