/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package samlv2

import "github.com/blend/go-sdk/ex"

// CanonicalizerAlgorithm identifies the XML canonicalization algorithm that the SAML provider should use to sign XML
type CanonicalizerAlgorithm string

// Supported canonicalization algorithms
const (
	CanonicalXML10ExclusiveAlgorithmID CanonicalizerAlgorithm = "http://www.w3.org/2001/10/xml-exc-c14n#"
	CanonicalXML11AlgorithmID          CanonicalizerAlgorithm = "http://www.w3.org/2006/12/xml-c14n11"
)

// error classes
const (
	ErrorUnsupportedCanonicalizer ex.Class = "Unsupported canonicalizer"
)
