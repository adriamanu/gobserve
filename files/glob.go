package files

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func isPatternAWildcard(pattern string) bool {
	if len(pattern) > 0 && pattern[:1] == "*" {
		return true
	}
	return false
}

func GlobFiles(pattern string) []string {
	tokenizedPattern := strings.Split(pattern, "/")
	patternLen := len(tokenizedPattern)

	var lookupPattern string
	// if the last token of the pattern contains a wildcard we lookup on the file extension
	// otherwise we lookup on the last token directly
	// a/aa/*.go -> we take .go as a lookup pattern because it contains a wildcar
	// a/aa/aaa.go -> we lookup for aaa.go file
	if isPatternAWildcard(tokenizedPattern[len(tokenizedPattern)-1]) {
		lookupPattern = "." + strings.SplitN(tokenizedPattern[len(tokenizedPattern)-1], ".", 2)[1]
	} else {
		lookupPattern = tokenizedPattern[len(tokenizedPattern)-1]
	}

	var fp []string
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Contains(path, lookupPattern) {
				fileInfo, err := os.Stat(path)
				if err != nil {
					log.Fatal(err)
				}
				if fileInfo.Mode().IsDir() {
					// exclude directories
					return nil
				}

				// current path tokenized
				cp := strings.Split(path, "/")
				cplen := len(cp)

				// files at root level
				if cplen == 1 {
					fp = append(fp, cp[0])
				} else if cplen <= patternLen {
					add := false
					for i := range cp {
						if (i == cplen-1) && strings.Contains(cp[i], lookupPattern) {
							add = true
							break
						}
						// break if dir pattern is not a double star or it does not match our defined pattern
						if !(tokenizedPattern[i] == "**") && !(cp[i] == tokenizedPattern[i]) {
							break
						}
					}
					if add {
						// if file matched our criteria we add it to files list
						fp = append(fp, path)
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
