/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

// BlockTypes
const (
	BlockTypeCertificate	= "CERTIFICATE"
	BlockTypeRSAPrivateKey	= "RSA PRIVATE KEY"
)

// Not After defaults.
const (
	DefaultCANotAfterYears		= 10
	DefaultClientNotAfterYears	= 1
	DefaultServerNotAfterYears	= 5
)
