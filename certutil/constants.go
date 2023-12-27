/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
