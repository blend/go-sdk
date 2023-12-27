/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package crypto

import "crypto/aes"

// Important constants.
const (
	// DefaultKeySize is the size of keys to generate for client use.
	DefaultKeySize	= 32
	// KeyVersionSize is the size of the key version prefix.
	KeyVersionSize	= (4 + 2 + 2 + 1)	// YYYY + MM + DD + :
	// IVSize is the size of the IV prefix.
	IVSize	= aes.BlockSize
	// HashSize is the size of the hash prefix.
	HashSize	= 32	// reasons.
)
