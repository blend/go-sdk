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
	"io"
	"strings"
)

// Source is a set of ownership entries.
type Source struct {
	Source	string
	Paths	[]Path
}

// WriteTo writes the owners to a given file.
func (s Source) WriteTo(wr io.Writer) (total int64, err error) {
	var n int
	n, err = fmt.Fprintf(wr, "%s%s\n", OwnersFileSourceComment, s.Source)
	if err != nil {
		return
	}
	total += int64(n)
	for _, entry := range s.Paths {
		n, err = fmt.Fprintf(wr, "%s %s\n", entry.PathGlob, strings.Join(entry.Owners, " "))
		if err != nil {
			return
		}
		total += int64(n)
	}
	n, err = fmt.Fprintf(wr, "%s%s\n", OwnersFileSourceEndComment, s.Source)
	total += int64(n)
	return
}
