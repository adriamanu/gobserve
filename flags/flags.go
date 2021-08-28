package flags

import (
	"flag"
	"fmt"
	"log"
)

var ConfigFlag = flag.String("config", "", "configuration file in json or yaml format")
var FilesFlag = flag.String("files", "", "files to watch")
var CommandFlag = flag.String("c", "", "command to execute")
var IgnoreFlag = flag.String("ignore", "", "regex of files to ignore")

func CheckFlags() {
	flag.Parse()

	errMsg := "Flag %s is mandatory"

	if *ConfigFlag == "" {
		// if config flag is not used we use other flags
		if *FilesFlag == "" {
			log.Fatal(fmt.Errorf(errMsg, "-files"))
		}
		if *CommandFlag == "" {
			log.Fatal(fmt.Errorf(errMsg, "-c"))
		}
	} else {
		// if config flag is used we disable other flags
		if *FilesFlag != "" {
			fmt.Println("-files flag has been ignored due to configuration file")
			*FilesFlag = ""
		}
		if *CommandFlag != "" {
			fmt.Println("-c flag has been ignored due to configuration file")
			*CommandFlag = ""
		}
		if *IgnoreFlag != "" {
			fmt.Println("-ignore flag has been ignored due to configuration file")
			*IgnoreFlag = ""
		}
	}
}
