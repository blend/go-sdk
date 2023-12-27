/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sh

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// PutOutput runs a given command with a given reader as its stdin and captures the output.
func PutOutput(stdin io.Reader, command string, args ...string) ([]byte, error) {
	absoluteCommand, err := exec.LookPath(command)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(absoluteCommand, args...)
	cmd.Env = os.Environ()
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	cmd.Stdin = stdin
	return buffer.Bytes(), cmd.Run()
}
