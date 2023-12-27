/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

// OutputParsed parses a command as binary and arguments.
// It resolves the command name in your $PATH list for you.
// It captures combined output and returns it as bytes.
func OutputParsed(statement string) ([]byte, error) {
	cmd, err := CmdParsed(statement)
	if err != nil {
		return nil, err
	}
	cmd.Env = os.Environ()
	return OutputCmd(cmd)
}

// OutputParsedContext parses a command as binary and arguments with a given context.
// It resolves the command name in your $PATH list for you.
// It captures combined output and returns it as bytes.
func OutputParsedContext(ctx context.Context, statement string) ([]byte, error) {
	cmd, err := CmdParsedContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	cmd.Env = os.Environ()
	return OutputCmd(cmd)
}

// Output runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It captures combined output and returns it as bytes.
func Output(command string, args ...string) ([]byte, error) {
	cmd, err := Cmd(command, args...)
	if err != nil {
		return nil, err
	}
	cmd.Env = os.Environ()
	return OutputCmd(cmd)
}

// OutputContext runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It captures combined output and returns it as bytes.
func OutputContext(ctx context.Context, command string, args ...string) ([]byte, error) {
	cmd, err := Cmd(command, args...)
	if err != nil {
		return nil, err
	}
	cmd.Env = os.Environ()
	return OutputCmd(cmd)
}

// OutputCmd captures the output of a given command.
func OutputCmd(cmd *exec.Cmd) ([]byte, error) {
	output := new(bytes.Buffer)
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Run(); err != nil {
		return output.Bytes(), err
	}
	return output.Bytes(), nil
}
