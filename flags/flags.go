package flags

import (
	"flag"
	"fmt"
	"log"
)

var ConfigFlag = flag.String("conf", "", "configuration file")
var IgnoreFlag = flag.String("ignore", "", "regex of files to ignore")
var FilesFlag = flag.String("files", "", "files to watch")
var CommandFlag = flag.String("c", "", "command to execute")

func CheckFlags() {
	flag.Parse()

	errMsg := "Flag %s is mandatory"
	if *FilesFlag == "" {
		log.Fatal(fmt.Errorf(errMsg, "-files"))
	}
	if *CommandFlag == "" {
		log.Fatal(fmt.Errorf(errMsg, "-c"))
	}
}
