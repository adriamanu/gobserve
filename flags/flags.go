package flags

import (
	"flag"
	"fmt"
	"log"
)

const config string = "config"
const files string = "files"
const ignore string = "ignore"
const command string = "c"

var ConfigFlag = flag.String(config, "", "configuration file in json or yaml format")
var FilesFlag = flag.String(files, "", "files to watch")
var CommandFlag = flag.String(command, "", "command to execute")
var IgnoreFlag = flag.String(ignore, "", "regex of files to ignore")

func CheckFlags() {
	flag.Parse()

	errMsg := "Flag %s is mandatory"

	if *ConfigFlag == "" {
		// if config flag is not used we use other flags
		if *FilesFlag == "" {
			log.Fatal(fmt.Errorf(errMsg, files))
		}
		if *CommandFlag == "" {
			log.Fatal(fmt.Errorf(errMsg, command))
		}
	} else {
		// if config flag is used we disable other flags
		excludeFlag(FilesFlag, files)
		excludeFlag(CommandFlag, command)
		excludeFlag(IgnoreFlag, ignore)
	}
}

func excludeFlag(f *string, flagName string) {
	if *f != "" {
		fmt.Printf("flag -%s has been ignored due to configuration file\n", flagName)
		*f = ""
	}
}
