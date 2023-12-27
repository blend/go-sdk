/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestContextWithTimestamp(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2019, 8, 16, 12, 11, 10, 9, time.UTC)
	assert.Equal(ts, GetTimestamp(WithTimestamp(context.Background(), ts)))
	assert.True(GetTimestamp(context.Background()).IsZero())
}

func TestContextWithRestricted(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, GetRestricted(WithRestricted(context.Background(), true)))
}

func TestContextWithPath(t *testing.T) {
	assert := assert.New(t)

	path := []string{"one", "two"}
	path2 := []string{"two", "three"}
	assert.Equal(path, GetPath(WithPath(context.Background(), path...)))
	assert.Equal(path, GetPath(WithPath(WithPath(context.Background(), path2...), path...)))
	assert.Nil(GetPath(context.Background()))
}

func TestContextWithSetLabels(t *testing.T) {
	assert := assert.New(t)

	labels := Labels{"one": "two"}
	labels2 := Labels{"two": "three"}
	assert.Equal(labels, GetLabels(WithSetLabels(context.Background(), labels)))
	assert.Equal(labels, GetLabels(WithSetLabels(WithSetLabels(context.Background(), labels2), labels)))
	assert.NotNil(GetLabels(context.Background()))
	assert.Empty(GetLabels(context.Background()))
}

func TestContextWithAnnotation(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	ctx = WithAnnotation(ctx, "one", "two")
	expectedAnnotations := Annotations{"one": "two"}
	assert.Equal(expectedAnnotations, GetAnnotations(ctx))

	ctx = WithAnnotation(ctx, "two", 3)
	expectedAnnotations = Annotations{
		"one":	"two",
		"two":	3,
	}
	assert.Equal(expectedAnnotations, GetAnnotations(ctx))
}

func TestContextWithLabels_Mutating(t *testing.T) {
	its := assert.New(t)

	ctx := context.Background()
	l0 := Labels{"one": "two", "three": "four"}
	ctx0 := WithLabels(ctx, l0)
	l1 := Labels{"one": "not-two", "two": "three"}
	ctx1 := WithLabels(ctx0, l1)

	l2 := GetLabels(ctx1)

	l2["foo"] = "bar"

	its.Equal("not-two", l2["one"])
	its.Equal("three", l2["two"])
	its.Equal("four", l2["three"])
	its.Equal("bar", l2["foo"])

	its.Equal("two", l0["one"])
	its.Empty(l0["foo"])
}

func TestContextWithLabel_Mutating(t *testing.T) {
	its := assert.New(t)

	original := Labels{"four": "five"}
	ctx := WithSetLabels(context.Background(), original)
	its.Equal("five", GetLabels(ctx)["four"])

	ctx0 := WithLabel(ctx, "one", "two")
	ctx1 := WithLabel(ctx0, "three", "four")
	ctx2 := WithLabel(ctx1, "four", "not-five")

	l2 := GetLabels(ctx2)

	its.Equal("two", l2["one"])
	its.Equal("four", l2["three"])
	its.Equal("not-five", l2["four"])

	its.Equal("five", original["four"])
}
