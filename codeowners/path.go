/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

// Path is a path in the codeowners file.
type Path struct {
	PathGlob	string
	Owners		[]string
}
