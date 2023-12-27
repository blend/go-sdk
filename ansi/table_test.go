/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ansi

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTable(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)

	err := Table(buffer, []string{"foo", "bar"}, [][]string{{"1", "2"}, {"3", "4"}})
	assert.Nil(err)
	assert.NotEmpty(buffer.String())
}

func TestTableEmptyErr(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)

	err := Table(buffer, nil, [][]string{{"1", "2"}, {"3", "4"}})
	assert.NotNil(err)
}

func TestTableWriteErr(t *testing.T) {
	assert := assert.New(t)

	failAfter := &failAfter{MaxBytes: 32}

	err := Table(failAfter, []string{"foo", "bar", "baz"}, [][]string{
		{"1", "2", "3"},
		{"4", "5", "6"},
		{"7", "8", "9"},
	})
	assert.NotNil(err)
	assert.Equal("did fail", err.Error())
}

type failAfter struct {
	Written		[]byte
	MaxBytes	int
}

func (fa *failAfter) Write(contents []byte) (int, error) {
	fa.Written = append(fa.Written, contents...)
	if len(fa.Written) > fa.MaxBytes {
		return len(contents), fmt.Errorf("did fail")
	}
	return len(contents), nil
}

func TestTableForSlice(t *testing.T) {
	assert := assert.New(t)

	objects := []struct {
		ID	int
		Name	string
	}{
		{1, "Foo"},
		{2, "Bar"},
		{3, "Baz"},
	}

	output := new(bytes.Buffer)
	assert.Nil(TableForSlice(output, objects))
	assert.Equal(
		"┌──┬────┐\n│ID│Name│\n├──┼────┤\n│1 │Foo │\n│2 │Bar │\n│3 │Baz │\n└──┴────┘\n",
		output.String(),
	)
}

func TestTableForSliceUnicode(t *testing.T) {
	assert := assert.New(t)

	objects := []struct {
		ID	string
		Count	int
	}{
		{"モ foo", 1},
		{"ふ bar", 1},
		{"ス baz", 3},
	}

	output := new(bytes.Buffer)
	assert.Nil(TableForSlice(output, objects))
	assert.Equal(
		"┌──────┬─────┐\n│ID    │Count│\n├──────┼─────┤\n│モ foo│1    │\n│ふ bar│1    │\n│ス baz│3    │\n└──────┴─────┘\n",
		output.String(),
	)
}
