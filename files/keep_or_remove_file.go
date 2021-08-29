package files

import (
	"goverwatch/helpers"
)

// ShouldKeepFile check if the provided file exists in the array of files passed as parameter.
// It returns true if the file doesn't exists yet in the array, otherwise it returns false.
func ShouldKeepFile(fileToKeep string, existingFiles []string) bool {
	for _, f := range existingFiles {
		if f == fileToKeep {
			return false
		}
	}
	return true
}

// RemoveFileFromList remove the provided file from the array of files.
func RemoveFileFromList(files *[]string, file string) {
	for i, f := range *files {
		if f == file {
			*files = helpers.Remove(*files, i)
		}
	}
	return
}

// Some files may match several patterns especially while using wildcard.
// RemoveDuplicatedFiles remove duplicated files from the two dimensional array passed as parameter and returns the "clean" array of files.
func RemoveDuplicatedFiles(f [][]string) []string {
	var filesToKeep []string
	if len(f) > 0 {
		filesToKeep = f[0]
	} else {
		return filesToKeep
	}

	for i := 1; i < len(f); i++ {
		for j := 0; j < len(f[i]); j++ {
			if ShouldKeepFile(f[i][j], filesToKeep) {
				filesToKeep = append(filesToKeep, f[i][j])
			}
		}
	}

	return filesToKeep
}

// RemoveIgnoredFiles removes ignored files from the array of files to keep.
func RemoveIgnoredFiles(filesToKeep *[]string, filesToIgnore []string) {
	for i := 0; i < len(filesToIgnore); i++ {
		if !ShouldKeepFile(filesToIgnore[i], *filesToKeep) {
			RemoveFileFromList(filesToKeep, filesToIgnore[i])
		}
	}
}
