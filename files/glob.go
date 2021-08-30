package files

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var tokenizedPattern []string
var patternLen int
var fp []string
var lookupPattern string

// init global variables
func initGlobalVariables(pattern string) {
	tokenizedPattern = strings.Split(pattern, "/")
	patternLen = len(tokenizedPattern)
	fp = []string{}

	// if the last token of the pattern contains a wildcard we lookup on the file extension
	// otherwise we lookup on the last token directly
	// a/aa/*.go -> we take .go as a lookup pattern because it contains a wildcar
	// a/aa/aaa.go -> we lookup for aaa.go file
	if isPatternAWildcard(tokenizedPattern[len(tokenizedPattern)-1]) {
		lookupPattern = "." + strings.SplitN(tokenizedPattern[len(tokenizedPattern)-1], ".", 2)[1]
	} else {
		lastToken := tokenizedPattern[len(tokenizedPattern)-1]
		// pattern may contains a single star
		if lastToken[:1] == "*" {
			lookupPattern = lastToken[1:]
		} else {
			lookupPattern = lastToken
		}
	}
}

// isPatternAWildcard returns a boolen if the pattern provided is a wildcard.
func isPatternAWildcard(pattern string) bool {
	if len(pattern) > 0 && (pattern[:2] == "**" || pattern[:2] == "*.") {
		return true
	}
	return false
}

func addMatchingFile(path string, info os.FileInfo, err error) error {
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
}

// GlobFiles returns all the files that match a pattern.
// * and ** patterns are handled.
func GlobFiles(pattern string) ([]string, error) {
	// init global variables
	initGlobalVariables(pattern)

	// walk through files and directories and add files that match to fp[] global variable
	err := filepath.Walk(".", addMatchingFile)

	if err != nil {
		return []string{}, nil
	}

	return fp, nil
}
