package sh

import (
	"bytes"
	"os"
	"os/exec"
)

// Output runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It captures combined output and returns it as bytes.
func Output(command string, args ...string) ([]byte, error) {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return nil, err
	}
	output := new(bytes.Buffer)
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
