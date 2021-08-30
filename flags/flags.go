package flags

import (
	"errors"
	"flag"
	"fmt"
)

const config string = "config"
const files string = "files"
const ignore string = "ignore"
const command string = "command"

var ConfigFlag = flag.String(config, "", "configuration file in json or yaml format")
var FilesFlag = flag.String(files, "", "files to watch")
var CommandFlag = flag.String(command, "", "command to execute")
var IgnoreFlag = flag.String(ignore, "", "regex of files to ignore")

// CheckFlags parse flags used with the cli and exit with an error message if a mandatory flag is missing.
func CheckFlags() error {
	flag.Parse()

	errMsg := "Flag %s is mandatory"

	if *ConfigFlag == "" {
		// if config flag is not used we use other flags
		if *FilesFlag == "" {
			return errors.New(fmt.Sprintf(errMsg, files))
		}
		if *CommandFlag == "" {
			return errors.New(fmt.Sprintf(errMsg, files))
		}
	} else {
		// if config flag is used we disable other flags
		disableFlag(FilesFlag, files)
		disableFlag(CommandFlag, command)
		disableFlag(IgnoreFlag, ignore)
	}
	return nil
}

// disableFlag disable flag when the flag -config is passed to the program.
func disableFlag(f *string, flagName string) {
	if *f != "" {
		fmt.Printf("flag -%s has been disabled due to configuration file\n", flagName)
		*f = ""
	}
}
