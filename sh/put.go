package sh

import (
	"io"
	"os"
	"os/exec"
)

// Put  runs a given command with a given reader as its stdin.
func Put(stdin io.Reader, command string, args ...string) error {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = stdin
	return cmd.Run()
}
