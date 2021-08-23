package files

import (
	"goverwatch/helpers"
)

func ShouldKeepFile(fileToKeep string, existingFiles []string) bool {
	for _, f := range existingFiles {
		if f == fileToKeep {
			return false
		}
	}
	return true
}

func RemoveFileFromList(files *[]string, file string) {
	for i, f := range *files {
		if f == file {
			*files = helpers.Remove(*files, i)
		}
	}
	return
}

func RemoveGlobDuplicates(f [][]string) []string {
	var filesToKeep []string
	if len(f) > 0 {
		filesToKeep = f[0]
	} else {
		return filesToKeep
	}

	// don't keep duplicates - some files may match several patterns
	for i := 1; i < len(f); i++ {
		for j := 0; j < len(f[i]); j++ {
			if ShouldKeepFile(f[i][j], filesToKeep) {
				filesToKeep = append(filesToKeep, f[i][j])
			}
		}
	}

	return filesToKeep
}

func RemoveIgnoredFiles(filesToKeep *[]string, filesToIgnore []string) {
	for i := 0; i < len(filesToIgnore); i++ {
		if !ShouldKeepFile(filesToIgnore[i], *filesToKeep) {
			RemoveFileFromList(filesToKeep, filesToIgnore[i])
		}
	}
}
