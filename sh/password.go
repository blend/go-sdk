package sh

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// MustPassword gives a prompt and reads input until newlines without printing the input to screen.
// The prompt is written to stdout with `fmt.Print` unchanged.
// It panics on error.
func MustPassword(prompt string) string {
	output, err := Password(prompt)
	if err != nil {
		panic(err)
	}
	return output
}

// Password prints a prompt and reads input until newlines without printing the input to screen.
// The prompt is written to stdout with `fmt.Print` unchanged.
func Password(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)
	results, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Fprintln(os.Stdout)
	return string(results), nil
}

// Passwordf gives a prompt and reads input until newlines without printing the input to screen.
// The prompt is written to stdout with `fmt.Printf` unchanged.
func Passwordf(format string, args ...interface{}) (string, error) {
	fmt.Fprintf(os.Stdout, format, args...)
	results, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Fprintln(os.Stdout)
	return string(results), nil
}
