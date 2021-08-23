package files

import (
	"testing"
)

var nestedPattern string = "**/**/*.go"
var simplePattern string = "*.go"

func TestRemoveFromList(t *testing.T) {
	files := GlobFiles(nestedPattern)
	t.Run("Remove file from list", func(t *testing.T) {
		amountOfFilesBeforeRemoving := len(files)
		filesToRemove := GlobFiles(simplePattern)
		for _, file := range filesToRemove {
			RemoveFileFromList(&files, file)
		}
		if !(len(files) == amountOfFilesBeforeRemoving-len(filesToRemove)) {
			t.Errorf("Amount of files should have been decreased by one.")
		}
	})

}

func TestRemoveIgnoreFiles(t *testing.T) {
	files := GlobFiles(nestedPattern)
	t.Run("Remove ignored files", func(t *testing.T) {
		amountOfFilesBeforeRemoving := len(files)
		filesToIgnore := GlobFiles("*_test.go")
		amountOfIgnoredFiles := len(filesToIgnore)
		RemoveIgnoredFiles(&files, filesToIgnore)
		if !(len(files) == amountOfFilesBeforeRemoving-amountOfIgnoredFiles) {
			t.Errorf("Amount of files should have been decreased by %d.", amountOfIgnoredFiles)
		}
	})
}

func TestRemoveDuplicatedFiles(t *testing.T) {
	allFiles := [][]string{}
	files := GlobFiles(nestedPattern)
	someMoreFiles := GlobFiles(simplePattern)
	allFiles = append(allFiles, files, someMoreFiles)

	cleanedFilesList := RemoveGlobDuplicates(allFiles)
	// files is a superset of somemorefiles as it contains the same files + some from other directories
	if len(cleanedFilesList) != len(files) {
		t.Errorf("Some files haven't been ignored.")
	}
}
