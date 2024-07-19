/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package uuid

// crypto/rand here for correctness.
import "crypto/rand"

// V4 Create a new UUID version 4.
func V4() UUID {
	var uuid UUID = Empty()
	_, _ = rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set variant 2
	return uuid
}
