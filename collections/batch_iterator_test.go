/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_BatchIterator(t *testing.T) {
	a := assert.New(t)

	bi := &BatchIterator[string]{BatchSize: 100}
	a.False(bi.HasNext())
	a.Empty(bi.Next())

	bi = &BatchIterator[string]{Items: generateBatchItems(10)}
	a.True(bi.HasNext())
	a.Empty(bi.Next())

	// handle edge case where somehow the cursor gets set beyond the
	// last element of the items.
	bi = &BatchIterator[string]{Items: generateBatchItems(10), Cursor: 15}
	a.False(bi.HasNext())
	a.Empty(bi.Next())

	bi = &BatchIterator[string]{Items: generateBatchItems(10), BatchSize: 100}
	a.True(bi.HasNext())
	a.Len(bi.Next(), 10)
	a.False(bi.HasNext())

	bi = &BatchIterator[string]{Items: generateBatchItems(100), BatchSize: 10}
	for x := 0; x < 10; x++ {
		a.True(bi.HasNext())
		a.Len(bi.Next(), 10, fmt.Sprintf("failed on pass %d", x))
	}
	a.False(bi.HasNext())

	bi = &BatchIterator[string]{Items: generateBatchItems(105), BatchSize: 10}
	for x := 0; x < 10; x++ {
		a.True(bi.HasNext())
		a.Len(bi.Next(), 10, fmt.Sprintf("failed on pass %d", x))
	}
	a.True(bi.HasNext())
	a.Len(bi.Next(), 5)
	a.False(bi.HasNext())
}

func generateBatchItems(count int) (output []string) {
	output = make([]string, count)
	for x := 0; x < count; x++ {
		output[x] = fmt.Sprint(x)
	}
	return output
}
