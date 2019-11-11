package sh

import (
	"os"
	"os/exec"
)

// Exec runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It does not show output.
func Exec(command string, args ...string) error {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
