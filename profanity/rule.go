/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

// Rule is a criteria for profanity.
type Rule interface {
	Check(file string, contents []byte) RuleResult
}
