/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package codeowners

// Path is a path in the codeowners file.
type Path struct {
	PathGlob string
	Owners   []string
}
