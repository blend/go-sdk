/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package diff

// Diff represents one diff operation
type Diff struct {
	Type Operation
	Text string
}
