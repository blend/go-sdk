/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sh

import (
	"context"
	"os"
)

// ExecParsed parses a command string into binary and args.
// It resolves the command name in your $PATH list for you.
// It does not show output.
func ExecParsed(command string) error {
	cmd, err := CmdParsed(command)
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	return cmd.Run()
}

// ExecParsedContext parses a command string into binary and args within a context.
// It resolves the command name in your $PATH list for you.
// It does not show output.
func ExecParsedContext(ctx context.Context, command string) error {
	cmd, err := CmdParsedContext(ctx, command)
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	return cmd.Run()
}

// Exec runs a command with a given list of arguments.
// It resolves the command name in your $PATH list for you.
// It does not show output.
func Exec(command string, args ...string) error {
	cmd, err := Cmd(command, args...)
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	return cmd.Run()
}

// ExecContext runs a command with a given list of arguments within a context.
// It resolves the command name in your $PATH list for you.
// It does not show output.
func ExecContext(ctx context.Context, command string, args ...string) error {
	cmd, err := CmdContext(ctx, command, args...)
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	return cmd.Run()
}
