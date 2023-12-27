/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import "github.com/blend/go-sdk/ex"

// Common error codes.
const (
	ErrNotFound			ex.Class	= "vault; not found"
	ErrUnauthorized			ex.Class	= "vault; not authorized"
	ErrServerError			ex.Class	= "vault; remote error"
	ErrBatchTransitEncryptError	ex.Class	= "vault; batch encryption error"
	ErrBatchTransitDecryptError	ex.Class	= "vault; batch decryption error"
)
