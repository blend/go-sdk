/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package diff

// Diff represents one diff operation
type Diff struct {
	Type	Operation
	Text	string
}
