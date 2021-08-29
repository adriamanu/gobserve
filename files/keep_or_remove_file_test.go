package files

import (
	"log"
	"testing"
)

var nestedPattern string = "**/**/*.go"
var simplePattern string = "*.go"

func TestRemoveFromList(t *testing.T) {
	files, err := GlobFiles(nestedPattern)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Remove file from list", func(t *testing.T) {
		amountOfFilesBeforeRemoving := len(files)

		filesToRemove, err := GlobFiles(simplePattern)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range filesToRemove {
			RemoveFileFromList(&files, file)
		}

		if !(len(files) == amountOfFilesBeforeRemoving-len(filesToRemove)) {
			t.Errorf("Amount of files should have been decreased by one.")
		}
	})
}

func TestRemoveIgnoreFiles(t *testing.T) {
	files, err := GlobFiles(nestedPattern)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Remove ignored files", func(t *testing.T) {
		amountOfFilesBeforeRemoving := len(files)

		filesToIgnore, err := GlobFiles("*_test.go")
		if err != nil {
			log.Fatal(err)
		}

		RemoveIgnoredFiles(&files, filesToIgnore)

		amountOfIgnoredFiles := len(filesToIgnore)
		if !(len(files) == amountOfFilesBeforeRemoving-amountOfIgnoredFiles) {
			t.Errorf("Amount of files should have been decreased by %d.", amountOfIgnoredFiles)
		}
	})
}

func TestRemoveDuplicatedFiles(t *testing.T) {
	allFiles := [][]string{}

	files, err := GlobFiles(nestedPattern)
	if err != nil {
		log.Fatal(err)
	}

	someMoreFiles, err := GlobFiles(simplePattern)
	if err != nil {
		log.Fatal(err)
	}

	allFiles = append(allFiles, files, someMoreFiles)

	cleanedFilesList := RemoveDuplicatedFiles(allFiles)

	// files is a superset of somemorefiles as it contains the same files + some from other directories
	if len(cleanedFilesList) != len(files) {
		t.Errorf("Some files haven't been ignored.")
	}
}
