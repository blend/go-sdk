/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

// Option is a tweak to selector parsing.
type Option func(p *Parser)

// SkipValidation is an option to skip checking the values of selector expressions.
func SkipValidation(p *Parser) {
	p.skipValidation = true
}
