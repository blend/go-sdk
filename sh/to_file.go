/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sh

import "os"

// MustToFile opens or creates a file and panics on error.
// This is meant to be used as an output writer for a command.
func MustToFile(path string) *os.File {
	file, err := ToFile(path)
	if err != nil {
		panic(err)
	}
	return file
}

// ToFile opens or creates a file.
// This is meant to be used as an output writer for a command.
func ToFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		return os.Open(path)
	}
	return os.Create(path)
}
