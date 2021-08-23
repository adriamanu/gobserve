package files

import (
	"fmt"
	"goverwatch/colors"
	"goverwatch/commands"
	"log"
	"os"
)

type fileToWatch struct {
	// location of the file
	filePath string
	// last time the file has been modified
	// if this value change the watcher will trigger a custom command
	modificationDate int64
}

// this slice if filled with globbed files
var filesToWatch []fileToWatch

// we glob a list of files to watch and pass them to this function
// it add to our slice
func DeclareFilesToWatch(files []string) {
	for i := range files {
		filesToWatch = append(filesToWatch, fileToWatch{files[i], getFileModificationDate(files[i])})
	}
}

// retrieve last modification date of a file
func getFileModificationDate(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.ModTime().UnixNano()
}

// compare the date with our original one
func fileHasBeenModified(originalModificationDate, modificationDate int64) bool {
	if originalModificationDate != modificationDate {
		return true
	}
	return false
}

func executeCommandOnFileModification(cmd commands.CommandToExecute, files []fileToWatch) {
	for _, file := range filesToWatch {
		originalModificationDate := file.modificationDate
		modificationDateToCheck := getFileModificationDate(file.filePath)
		if fileHasBeenModified(originalModificationDate, modificationDateToCheck) {
			fmt.Printf(colors.Yellow+"File %s has been modified\n"+colors.Reset, file.filePath)
			file.modificationDate = modificationDateToCheck
			commands.ExecCmd(cmd)
		}
	}
}

func Watch(cmd commands.CommandToExecute) {
	for {
		executeCommandOnFileModification(cmd, filesToWatch)
	}
}
