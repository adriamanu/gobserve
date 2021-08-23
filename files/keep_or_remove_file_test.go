package files

import (
	"testing"
)

var pattern string = "**/**/*.go"

func TestRemoveFromList(t *testing.T) {
	files := GlobFiles(pattern)
	t.Run("Remove file from list", func(t *testing.T) {
		amountOfFilesBeforeRemoving := len(files)
		filesToRemove := GlobFiles("*.go")
		for _, file := range filesToRemove {
			RemoveFileFromList(&files, file)
		}
		if !(len(files) == amountOfFilesBeforeRemoving-len(filesToRemove)) {
			t.Errorf("amount of files should have been decreased by one.")
		}
	})

}

func TestRemoveIgnoreFiles(t *testing.T) {
	files := GlobFiles(pattern)
	t.Run("Remove ignored files", func(t *testing.T) {
		amountOfFilesBeforeRemoving := len(files)
		filesToIgnore := GlobFiles("*_test.go")
		amountOfIgnoredFiles := len(filesToIgnore)
		RemoveIgnoredFiles(&files, filesToIgnore)
		if !(len(files) == amountOfFilesBeforeRemoving-amountOfIgnoredFiles) {
			t.Errorf("amount of files should have been decreased by %d.", amountOfIgnoredFiles)
		}
	})
}
