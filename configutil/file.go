/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"os"
	"strings"
)

// File reads the string contents of a file as a literal config value.
type File string

// String returns the string contents of a file.
func (f File) String(ctx context.Context) (*string, error) {
	contents, err := os.ReadFile(string(f))
	if err != nil {
		return nil, nil
	}
	stringContents := strings.TrimSpace(string(contents))
	return &stringContents, nil
}
