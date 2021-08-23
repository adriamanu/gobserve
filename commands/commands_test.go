package commands

import (
	"goverwatch/process"
	"testing"
)

func TestProcessAndCommands(t *testing.T) {
	var cmd CommandToExecute = ParseCmd("go list")
	ExecCmd(cmd)

	t.Run("Execute generated command", func(t *testing.T) {
		if process.RunningProcess == nil {
			t.Errorf("Running process has not been initialized, command has not been executed.")
		}
	})

	t.Run("Re execute command", func(t *testing.T) {
		if process.RunningProcess == nil {
			t.Errorf("Running process has not been initialized, command has not been executed.")
		}
		previousPid := process.RunningProcess.Pid
		ExecCmd(cmd)
		newPid := process.RunningProcess.Pid
		if previousPid == newPid {
			t.Errorf("Old and new pid can't be the same, it means that the process hasn't been killed and renewed.")
		}
	})

}
