/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import "io"

// File is a collection of codeowner sources.
type File []Source

// WriteTo writes the file to a given writer.
func (f File) WriteTo(wr io.Writer) (total int64, err error) {
	var n int
	var ln int64
	n, err = io.WriteString(wr, OwnersFileHeader+"\n")
	if err != nil {
		return
	}
	total += int64(n)

	for _, co := range f {
		ln, err = co.WriteTo(wr)
		if err != nil {
			return
		}
		total += int64(ln)
	}
	return
}
