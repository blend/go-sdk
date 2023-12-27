/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import "strings"

// ValidationErrors is a set of errors.
type ValidationErrors []error

// Error implements error.
func (ve ValidationErrors) Error() string {
	var output []string
	for _, e := range ve {
		output = append(output, e.Error())
	}
	return strings.Join(output, "\n")
}
