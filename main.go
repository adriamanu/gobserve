package main

import (
	"fmt"
	"goverwatch/colors"
	"goverwatch/commands"
	"goverwatch/files"
	"goverwatch/flags"
	"goverwatch/process"
	"strings"
)

func main() {
	process.CatchSigTerm()
	flags.CheckFlags()

	fmt.Printf(colors.Yellow + "And now my watch begins. It shall not end until my death.\n\n" + colors.Reset)

	// allow multiple patterns separated by a space
	patternsToGlob := strings.Split(*flags.FilesFlag, " ")
	var filesToGlob [][]string
	for i := range patternsToGlob {
		p := patternsToGlob[i]
		if p != "" {
			filesToGlob = append(filesToGlob, files.GlobFiles(p))
		}
	}

	var filesToIgnore [][]string
	patternsToIgnore := strings.Split(*flags.IgnoreFlag, " ")
	for j := range patternsToIgnore {
		pi := patternsToIgnore[j]
		if pi != "" {
			filesToIgnore = append(filesToIgnore, files.GlobFiles((pi)))
		}
	}

	f := files.RemoveGlobDuplicates(filesToGlob)
	ignoredFiles := files.RemoveGlobDuplicates(filesToIgnore)
	files.RemoveIgnoredFiles(&f, ignoredFiles)
	files.DeclareFilesToWatch(f)
	fmt.Printf(colors.Yellow+"watching on %s\n\n"+colors.Reset, f)

	cmd := commands.ParseCmd(*flags.CommandFlag)

	commands.ExecCmd(cmd)
	files.Watch(cmd)
}
