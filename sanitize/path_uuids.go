/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sanitize

import (
	"strings"

	"github.com/blend/go-sdk/uuid"
)

var (
	_ PathSanitizerFunc = PathUUIDs
)

// PathUUIDs is a path sanitizer func
// that replaces any uuids in a path with "?".
func PathUUIDs(path string) string {
	if path == "" || path == "/" {
		return path
	}

	pathParts := strings.Split(path, "/")
	for index := range pathParts {
		if id, _ := uuid.Parse(pathParts[index]); !id.IsZero() {
			pathParts[index] = "?"
		}
	}
	return strings.Join(pathParts, "/")
}
