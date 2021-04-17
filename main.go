package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

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

func watchForModification(filesToWatch []fileToWatch) {
	for {
		for i := range filesToWatch {
			originalModificationDate := filesToWatch[i].modificationDate
			modificationDateToCheck := getFileModificationDate(filesToWatch[i].filePath)
			if fileHasBeenModified(originalModificationDate, modificationDateToCheck) {
				fmt.Printf("File %s has been modified\n", filesToWatch[i].filePath)
				filesToWatch[i].modificationDate = modificationDateToCheck
				fmt.Println("**** RELOAD ****")
			}
		}
		time.Sleep(100 * time.Millisecond) // pause for 100 ms
	}
}

func main() {
	declareFilesToObserve([]string{"test.txt", "test2.txt"})
	watchForModification(filesToObserve)
}
