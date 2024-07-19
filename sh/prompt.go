/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/blend/go-sdk/ex"
)

// ErrUnexpectedNewLine is returned from scan.go when you just hit enter with nothing in the prompt
const ErrUnexpectedNewLine ex.Class = "unexpected newline"

// Prompt gives a prompt and reads input until newlines.
func Prompt(prompt string) string {
	return PromptFrom(os.Stdout, os.Stdin, prompt)
}

// Promptf gives a prompt of a given format and args and reads input until newlines.
func Promptf(format string, args ...interface{}) string {
	return PromptFrom(os.Stdout, os.Stdin, fmt.Sprintf(format, args...))
}

// PromptFrom gives a prompt and reads input until newlines from a given set of streams.
func PromptFrom(stdout io.Writer, stdin io.Reader, prompt string) string {
	fmt.Fprint(stdout, prompt)

	scanner := bufio.NewScanner(stdin)
	var output string
	if scanner.Scan() {
		output = scanner.Text()
	}
	return output
}
