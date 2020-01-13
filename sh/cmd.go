package sh

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"github.com/blend/go-sdk/ex"
)

// MustCmds returns a list of commands for a given set of statements.
// It is useful for running a batch of commands.
// It panics on error.
func MustCmds(statements ...string) []*exec.Cmd {
	var output []*exec.Cmd
	for _, statement := range statements {
		output = append(output, MustCmdParsed(statement))
	}
	return output
}

// Cmds returns a list of commands for a given set of statements.
// It is useful for running a batch of commands.
// It will return the first error it encounters.
func Cmds(statements ...string) ([]*exec.Cmd, error) {
	var output []*exec.Cmd
	for _, statement := range statements {
		cmd, err := CmdParsed(statement)
		if err != nil {
			return nil, err
		}
		output = append(output, cmd)
	}
	return output, nil
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
	parts := strings.Split(statement, " ")
	if len(parts) > 1 {
		return Cmd(parts[0], parts[1:]...)
	}
	return Cmd(parts[0])
}

// CmdParsedContext returns a command for a full comamnd statement within a context..
func CmdParsedContext(ctx context.Context, statement string) (*exec.Cmd, error) {
	parts := strings.SplitN(statement, " ", 2)
	if len(parts) > 1 {
		return Cmd(parts[0], parts[1])
	}
	return CmdContext(ctx, parts[0])
}

// MustCmd returns a new command with the fully qualified path of the executable.
// It panics on error.
func MustCmd(command string, args ...string) *exec.Cmd {
	cmd, err := Cmd(command, args...)
	if err != nil {
		panic(ex.New(err, ex.OptMessagef("looking for: %s", command)))
	}
	return cmd
}

// Cmd returns a new command with the fully qualified path of the executable.
func Cmd(command string, args ...string) (*exec.Cmd, error) {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return nil, ex.New(err, ex.OptMessagef("looking for: %s", command))
	}
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	return cmd, nil
}

// CmdContext returns a new command with the fully qualified path of the executable within a context.
func CmdContext(ctx context.Context, command string, args ...string) (*exec.Cmd, error) {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return nil, ex.New(err, ex.OptMessagef("looking for: %s", command))
	}
	cmd := exec.CommandContext(ctx, absoluteCommand, args...)
	cmd.Env = os.Environ()
	return cmd, nil
}
