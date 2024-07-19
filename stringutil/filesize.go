/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import "fmt"

// FileSize returns a string representation of a file size in bytes.
func FileSize(sizeBytes int) string {
	if sizeBytes >= 1<<30 {
		return fmt.Sprintf("%dgB", sizeBytes/(1<<30))
	} else if sizeBytes >= 1<<20 {
		return fmt.Sprintf("%dmB", sizeBytes/(1<<20))
	} else if sizeBytes >= 1<<10 {
		return fmt.Sprintf("%dkB", sizeBytes/(1<<10))
	}
	return fmt.Sprintf("%dB", sizeBytes)
}
