/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reflectutil

type testType struct {
	ID		int
	Name		string
	NotTagged	string
	Tagged		string
	SubTypes	[]subType
}

type subType struct {
	ID	int
	Name	string
}
