package files

import (
	"goverwatch/commands"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	t.Run("Watch over a list of files", func(t *testing.T) {
		// init filesToWatch
		DeclareFilesToWatch(GlobFiles("*.go"))
		// loop over list of files and execute command if files have been updated
		executeCommandOnFileModification(commands.ParseCmd("go list"), filesToWatch)
		// update file
		filesToWatch[0].modificationDate = time.Now().UnixNano()
		// re execute loop
		executeCommandOnFileModification(commands.ParseCmd("go list"), filesToWatch)
	})
}
