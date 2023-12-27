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
