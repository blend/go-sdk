package sh

import (
	"fmt"
	"os"
)

// Fatal exits the process with a given error.
func Fatal(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Fatalf exits the process with a given formatted message to stderr.
func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
