package sh

import (
	"os"
	"os/exec"
	"strings"
)

// MustCmds returns a command for a full comamnd statement.
func MustCmds(statements ...string) []*exec.Cmd {
	var output []*exec.Cmd
	for _, statement := range statements {
		output = append(output, MustCmdParsed(statement))
	}
	return output
}

// MustCmdParsed returns a command for a full comamnd statement.
func MustCmdParsed(statement string) *exec.Cmd {
	cmd, err := CmdParsed(statement)
	if err != nil {
		panic(err)
	}
	return cmd
}

// CmdParsed returns a command for a full comamnd statement.
func CmdParsed(statement string) (*exec.Cmd, error) {
	parts := strings.SplitN(statement, " ", 2)
	if len(parts) > 1 {
		return Cmd(parts[0], parts[1])
	}
	return Cmd(parts[0])
}

// MustCmd returns a new command with the fully qualified path of the executable.
// It panics on error.
func MustCmd(command string, args ...string) *exec.Cmd {
	cmd, err := Cmd(command, args...)
	if err != nil {
		panic(err)
	}
	return cmd
}

// Cmd returns a new command with the fully qualified path of the executable.
func Cmd(command string, args ...string) (*exec.Cmd, error) {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	return cmd, nil
}
