/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

// Option is a tweak to selector parsing.
type Option func(p *Parser)

// SkipValidation is an option to skip checking the values of selector expressions.
func SkipValidation(p *Parser) {
	p.skipValidation = true
}
