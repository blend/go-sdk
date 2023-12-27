/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import (
	"fmt"
	"strings"
)

// ParsePath parses a path line into a path and owners.
func ParsePath(pathLine string) (output Path, err error) {
	parts := strings.Split(pathLine, " ")
	if len(parts) < 2 {
		err = fmt.Errorf("invalid codeowners path line: %q", pathLine)
		return
	}
	output.PathGlob = parts[0]
	for _, owner := range parts[1:] {
		owner = strings.TrimSpace(owner)
		if owner != "" {
			output.Owners = append(output.Owners, owner)
		}
	}
	return
}
