/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

import (
	"context"
	"fmt"
	"os"
)

// RemoveFile removes a file and prints a debug message if the context sets that flag.
func RemoveFile(ctx context.Context, path string) error {
	if info, err := os.Stat(path); err != nil {
		return err
	} else if info.IsDir() {
		return fmt.Errorf("cannot remove file; %s is a directory", path)
	}
	return os.Remove(path)
}
