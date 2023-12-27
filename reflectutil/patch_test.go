/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reflectutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPatch(t *testing.T) {
	assert := assert.New(t)

	myObj := testType{}
	myObj.ID = 123
	myObj.Name = "Test Object"
	myObj.NotTagged = "Not Tagged"
	myObj.Tagged = "Is Tagged"
	myObj.SubTypes = append([]subType{}, subType{1, "One"})
	myObj.SubTypes = append(myObj.SubTypes, subType{2, "Two"})
	myObj.SubTypes = append(myObj.SubTypes, subType{3, "Three"})
	myObj.SubTypes = append(myObj.SubTypes, subType{4, "Four"})

	patchData := make(map[string]interface{})
	patchData["Tagged"] = "Is Not Tagged"

	err := Patch(&myObj, patchData)
	assert.Nil(err)
	assert.Equal("Is Not Tagged", myObj.Tagged)
}
