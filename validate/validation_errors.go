/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
