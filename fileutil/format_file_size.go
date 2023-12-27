/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil

import "strconv"

const (
	// Kilobyte represents the bytes in a kilobyte.
	Kilobyte int64 = 1 << 10
	// Megabyte represents the bytes in a megabyte.
	Megabyte int64 = Kilobyte << 10
	// Gigabyte represents the bytes in a gigabyte.
	Gigabyte int64 = Megabyte << 10
	// Terabyte represents the bytes in a terabyte.
	Terabyte int64 = Gigabyte << 10
)

// FormatFileSize returns a string representation of a file size in bytes.
func FormatFileSize(sizeBytes int64) string {
	switch {
	case sizeBytes >= 1<<40:
		return strconv.FormatInt(sizeBytes/Terabyte, 10) + "tb"
	case sizeBytes >= 1<<30:
		return strconv.FormatInt(sizeBytes/Gigabyte, 10) + "gb"
	case sizeBytes >= 1<<20:
		return strconv.FormatInt(sizeBytes/Megabyte, 10) + "mb"
	case sizeBytes >= 1<<10:
		return strconv.FormatInt(sizeBytes/Kilobyte, 10) + "kb"
	}
	return strconv.FormatInt(sizeBytes, 10)
}
