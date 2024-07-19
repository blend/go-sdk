/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

// Rule is a criteria for profanity.
type Rule interface {
	Check(file string, contents []byte) RuleResult
}
