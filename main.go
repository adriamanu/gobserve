package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// STRUCTURES
type fileToWatch struct {
	// location of the file
	filePath string
	// last time the file has been modified
	// if this value change the watcher will trigger a custom command
	modificationDate int64
}

type commandToExecute struct {
	command string
	args    []string
}

// VARS
var Reset = "\033[0m"
var Yellow = "\033[33m"

// this slice if filled with globbed files
var filesToWatch []fileToWatch

// currently running sub process
var runningProcess *os.Process

// FLAGS
var configFlag = flag.String("conf", "", "configuration file")
var ignoreFlag = flag.String("ignore", "", "regex of files to ignore")
var filesFlag = flag.String("files", "", "files to watch")
var commandFlag = flag.String("c", "", "command to execute")

func checkFlags() {
	errMsg := "Flag %s is mandatory"
	if *filesFlag == "" {
		log.Fatal(fmt.Errorf(errMsg, "-files"))
	}
	if *commandFlag == "" {
		log.Fatal(fmt.Errorf(errMsg, "-c"))
	}
}

// PROCESS MANAGEMENT
func killRunningProcess() {
	if runningProcess != nil {
		runningProcess.Signal(syscall.SIGKILL)
		// by using minus operator Kill will send signal to process group id
		// by doing so it will also kill process childrens
		err := syscall.Kill(-runningProcess.Pid, syscall.SIGKILL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func catchSigTerm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Printf(Yellow + "\n\nAnd now my watch is ended.\n" + Reset)
		killRunningProcess()
		os.Exit(0)
	}()
}

// COMMANDS
func parseCmd(commandToParse string) commandToExecute {
	tokenizedCommand := strings.Split(commandToParse, " ")
	amountOfToken := len(tokenizedCommand)
	cmd := tokenizedCommand[0]

	args := []string{}
	if amountOfToken > 1 {
		for i := 1; i < amountOfToken; i++ {
			args = append(args, tokenizedCommand[i])
		}
	}
	return commandToExecute{cmd, args}
}

func execCmd(c commandToExecute) {
	killRunningProcess()

	cmd := exec.Command(c.command, c.args...)
	// set the same pgid on childrens
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// redirect command errors and output to standard output and errors
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	runningProcess = cmd.Process
}

// FILES
func isPatternAWildcard(pattern string) bool {
	if len(pattern) > 0 && pattern[:1] == "*" {
		return true
	}
	return false
}

func globFiles(pattern string) []string {
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

// we glob a list of files to watch and pass them to this function
// it add to our slice
func declareFilesToWatch(files []string) {
	for i := range files {
		filesToWatch = append(filesToWatch, fileToWatch{files[i], getFileModificationDate(files[i])})
	}
}

func getFileModificationDate(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.ModTime().UnixNano()
}

func fileHasBeenModified(originalModificationDate, modificationDate int64) bool {
	if originalModificationDate != modificationDate {
		return true
	}
	return false
}

func watch(filesToWatch []fileToWatch, cmd commandToExecute) {
	for {
		for i := range filesToWatch {
			originalModificationDate := filesToWatch[i].modificationDate
			modificationDateToCheck := getFileModificationDate(filesToWatch[i].filePath)
			if fileHasBeenModified(originalModificationDate, modificationDateToCheck) {
				fmt.Printf(Yellow+"File %s has been modified\n"+Reset, filesToWatch[i].filePath)
				filesToWatch[i].modificationDate = modificationDateToCheck
				execCmd(cmd)
			}
		}
	}
}

func shouldKeepFile(fileToKeep string, existingFiles []string) bool {
	for _, f := range existingFiles {
		if f == fileToKeep {
			return false
		}
	}
	return true
}

func removeGlobDuplicates(f [][]string) []string {
	var filesToKeep []string
	if len(f) > 0 {
		filesToKeep = f[0]
	} else {
		return filesToKeep
	}

	// don't keep duplicates - some files may match several patterns
	for i := 1; i < len(f); i++ {
		for j := 0; j < len(f[i]); j++ {
			if shouldKeepFile(f[i][j], filesToKeep) {
				filesToKeep = append(filesToKeep, f[i][j])
			}
		}
	}

	return filesToKeep
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeFileFromList(files *[]string, file string) {
	for i, f := range *files {
		if f == file {
			*files = remove(*files, i)
		}
	}
	return
}

func removeIgnoredFiles(filesToKeep *[]string, filesToIgnore []string) {
	for i := 0; i < len(filesToIgnore); i++ {
		if !shouldKeepFile(filesToIgnore[i], *filesToKeep) {
			removeFileFromList(filesToKeep, filesToIgnore[i])
		}
	}
}

func main() {
	catchSigTerm()
	flag.Parse()
	checkFlags()

	fmt.Printf(Yellow + "And now my watch begins. It shall not end until my death.\n\n" + Reset)

	// allow multiple patterns separated by a space
	patternsToGlob := strings.Split(*filesFlag, " ")
	var filesToGlob [][]string
	for i := range patternsToGlob {
		p := patternsToGlob[i]
		if p != "" {
			filesToGlob = append(filesToGlob, globFiles(p))
		}
	}

	var filesToIgnore [][]string
	patternsToIgnore := strings.Split(*ignoreFlag, " ")
	for j := range patternsToIgnore {
		pi := patternsToIgnore[j]
		if pi != "" {
			filesToIgnore = append(filesToIgnore, globFiles((pi)))
		}
	}

	files := removeGlobDuplicates(filesToGlob)
	ignoredFiles := removeGlobDuplicates(filesToIgnore)
	removeIgnoredFiles(&files, ignoredFiles)
	declareFilesToWatch(files)
	fmt.Printf(Yellow+"watching on %s\n\n"+Reset, files)

	cmd := parseCmd(*commandFlag)

	execCmd(cmd)
	watch(filesToWatch, cmd)
}
