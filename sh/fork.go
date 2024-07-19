/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"context"
	"os"
	"os/exec"
)

// ForkParsed parses a command into binary and arguments.
// It resolves the command name in your $PATH list for you.
// It shows output and allows input.
// It is intended to be used in scripting.
func ForkParsed(command string) error {
	cmd, err := CmdParsed(command)
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// ForkParsedContext parses a command into binary and arguments.
// It resolves the command name in your $PATH list for you.
// It shows output and allows input.
// It is intended to be used in scripting.
func ForkParsedContext(ctx context.Context, command string) error {
	cmd, err := CmdParsedContext(ctx, command)
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Fork runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It shows output and allows input.
// It is intended to be used in scripting.
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
	return cmd.Run()
}

// ForkContext runs a command with a given list of arguments and context.
// It resolves the command name in your $PATH list for you.
// It shows output and allows input.
// It is intended to be used in scripting.
func ForkContext(ctx context.Context, command string, args ...string) error {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, absoluteCommand, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
