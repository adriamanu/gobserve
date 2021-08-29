package main

import (
	"fmt"
	"goverwatch/colors"
	"goverwatch/commands"
	"goverwatch/config"
	"goverwatch/files"
	"goverwatch/flags"
	"goverwatch/process"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	process.CatchSignalsAndExit()
	err := flags.CheckFlags()
	if err != nil {
		log.Fatal(err)
	}

	var patternsToGlob []string
	var patternsToIgnore []string
	var command string

	if *flags.ConfigFlag != "" {
		absolutePath, err := filepath.Abs(*flags.ConfigFlag)
		if err != nil {
			log.Fatal(err)
		}

		conf, err := config.ParseConfigFile(absolutePath)
		if err != nil {
			log.Fatal(err)
		}

		patternsToGlob = conf.Files
		patternsToIgnore = conf.IgnoredFiles
		command = conf.Command
	} else {
		patternsToGlob = strings.Split(*flags.FilesFlag, " ")
		patternsToIgnore = strings.Split(*flags.IgnoreFlag, " ")
		command = *flags.CommandFlag
	}

	var filesToGlob [][]string
	for i := range patternsToGlob {
		pattern := patternsToGlob[i]
		if pattern != "" {
			f, err := files.GlobFiles(pattern)
			if err != nil {
				log.Fatal(err)
			}
			filesToGlob = append(filesToGlob, f)
		}
	}

	var filesToIgnore [][]string
	for j := range patternsToIgnore {
		ignoredPattern := patternsToIgnore[j]
		if ignoredPattern != "" {
			f, err := files.GlobFiles(ignoredPattern)
			if err != nil {
				log.Fatal(err)
			}
			filesToIgnore = append(filesToIgnore, f)
		}
	}

	f := files.RemoveDuplicatedFiles(filesToGlob)
	ignoredFiles := files.RemoveDuplicatedFiles(filesToIgnore)
	files.RemoveIgnoredFiles(&f, ignoredFiles)

	fmt.Printf(colors.Yellow + "And now my watch begins. It shall not end until my death.\n\n" + colors.Reset)
	files.DeclareFilesToWatch(f)
	fmt.Printf(colors.Yellow+"watching on %s\n\n"+colors.Reset, f)

	cmd := commands.ParseCmd(command)
	commands.ExecCmd(cmd)
	files.Watch(cmd)
}
