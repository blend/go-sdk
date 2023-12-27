/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package shardutil

import "hash/fnv"

// Hash hashes a given string as an integer.
func Hash(value []byte) int {
	h := fnv.New32a()
	_, _ = h.Write(value)
	return int(h.Sum32())
}

// HashString hashes a given string as an integer.
func HashString(value string) int {
	return Hash([]byte(value))
}
