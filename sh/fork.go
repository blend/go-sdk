package sh

import (
	"os"
	"os/exec"
)

// Fork runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It shows output and allows input.
func Fork(command string, args ...string) error {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
