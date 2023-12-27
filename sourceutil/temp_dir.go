/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

import (
	"fmt"
	"os"
	"time"
)

// CreateTempDir creates a temporary directory with a given prefix.
func CreateTempDir(prefix string) (*TempDir, error) {
	fullPath := fmt.Sprintf("%s_%d", prefix, time.Now().UTC().UnixNano())
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return nil, err
	}
	return &TempDir{fullPath}, nil
}

// TempDir is a directory that can be cleaned up with Close.
type TempDir struct {
	Path string
}

// Close removes the directory.
func (td TempDir) Close() error {
	return os.RemoveAll(td.Path)
}
