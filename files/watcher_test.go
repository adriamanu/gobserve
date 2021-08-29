package files

import (
	"goverwatch/commands"
	"log"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	t.Run("Watch over a list of files", func(t *testing.T) {
		pattern := "*.go"

		// init filesToWatch
		files, err := GlobFiles(pattern)
		if err != nil {
			log.Fatal(err)
		}

		DeclareFilesToWatch(files)

		// loop over list of files and execute command if files have been updated
		executeCommandOnFileModification(commands.ParseCmd("go list"), filesToWatch)

		// update file
		filesToWatch[0].modificationDate = time.Now().UnixNano()

		// re execute loop
		executeCommandOnFileModification(commands.ParseCmd("go list"), filesToWatch)
	})
}
