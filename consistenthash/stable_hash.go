/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package consistenthash

import (
	"hash/crc64"
)

var (
	// stableCRC implements a stable crc64 hash table.
	// this allows us to have a consistent hash assignment
	// between process restarts.
	stableCRC = crc64.MakeTable(0xC96C5795D7870F42)
)

// HashFunction is a function that can be used to hash items.
type HashFunction func([]byte) uint64

// StableHash implements the default hash function with
// a stable crc64 table checksum.
func StableHash(data []byte) uint64 {
	return crc64.Checksum(data, stableCRC)
}
