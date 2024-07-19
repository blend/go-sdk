/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
