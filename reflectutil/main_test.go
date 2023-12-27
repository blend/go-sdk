/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reflectutil

type testType struct {
	ID        int
	Name      string
	NotTagged string
	Tagged    string
	SubTypes  []subType
}

type subType struct {
	ID   int
	Name string
}
