/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package semver

import (
	"fmt"
	"sort"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewVersion(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		version string
		err     bool
	}{
		{"1.2.3", false},
		{"1.0", false},
		{"1", false},
		{"1.2.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hypen", false},
		{"1.2.3-rc1-with-hypen", false},
		{"1.2.3.4", false},
		{"1.2.0.4-x.Y.0+metadata", false},
		{"1.2.0.4-x.Y.0+metadata-width-hypen", false},
		{"1.2.0-X-1.2.0+metadata~dist", false},
		{"1.2.3.4-rc1-with-hypen", false},
		{"1.2.3.4", false},
		{"v1.2.3", false},
		{"foo1.2.3", true},
		{"1.7rc2", false},
		{"v1.7rc2", false},
	}

	for _, tc := range cases {
		_, err := NewVersion(tc.version)
		assert.False(tc.err && err == nil, fmt.Sprintf("expected error for version: %s", tc.version))
		assert.False(!tc.err && err != nil, fmt.Sprintf("error for version %s: %s", tc.version, err))
	}
}

func TestVersionCompare(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.4.5", -1},
		{"1.2-beta", "1.2-beta", 0},
		{"1.2", "1.1.4", 1},
		{"1.2", "1.2-beta", 1},
		{"1.2+foo", "1.2+beta", 0},
		{"v1.2", "v1.2-beta", 1},
		{"v1.2+foo", "v1.2+beta", 0},
		{"v1.2.3.4", "v1.2.3.4", 0},
		{"v1.2.0.0", "v1.2", 0},
		{"v1.2.0.0.1", "v1.2", 1},
		{"v1.2", "v1.2.0.0", 0},
		{"v1.2", "v1.2.0.0.1", -1},
		{"v1.2.0.0", "v1.2.0.0.1", -1},
		{"v1.2.3.0", "v1.2.3.4", -1},
		{"1.7rc2", "1.7rc1", 1},
		{"1.7rc2", "1.7", -1},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", 1},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		assert.Nil(err)

		v2, err := NewVersion(tc.v2)
		assert.Nil(err)

		actual := v1.Compare(v2)
		expected := tc.expected
		assert.Equal(expected, actual, fmt.Sprintf("%s <=> %s\nexpected: %d\nactual: %d", tc.v1, tc.v2, expected, actual))
	}
}

func TestComparePreReleases(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2-beta.2", "1.2-beta.2", 0},
		{"1.2-beta.1", "1.2-beta.2", -1},
		{"1.2-beta.2", "1.2-beta.11", -1},
		{"3.2-alpha.1", "3.2-alpha", 1},
		{"1.2-beta.2", "1.2-beta.1", 1},
		{"1.2-beta.11", "1.2-beta.2", 1},
		{"1.2-beta", "1.2-beta.3", -1},
		{"1.2-alpha", "1.2-beta.3", -1},
		{"1.2-beta", "1.2-alpha.3", 1},
		{"3.0-alpha.3", "3.0-rc.1", -1},
		{"3.0-alpha3", "3.0-rc1", -1},
		{"3.0-alpha.1", "3.0-alpha.beta", -1},
		{"5.4-alpha", "5.4-alpha.beta", 1},
		{"v1.2-beta.2", "v1.2-beta.2", 0},
		{"v1.2-beta.1", "v1.2-beta.2", -1},
		{"v3.2-alpha.1", "v3.2-alpha", 1},
		{"v3.2-rc.1-1-g123", "v3.2-rc.2", 1},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		assert.Nil(err)

		v2, err := NewVersion(tc.v2)
		assert.Nil(err)

		actual := v1.Compare(v2)
		expected := tc.expected
		assert.Equal(expected, actual)
	}
}

func TestVersionMetadata(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		version  string
		expected string
	}{
		{"1.2.3", ""},
		{"1.2-beta", ""},
		{"1.2.0-x.Y.0", ""},
		{"1.2.0-x.Y.0+metadata", "metadata"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "metadata~dist"},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		assert.Nil(err)

		actual := v.Metadata()
		expected := tc.expected
		assert.Equal(expected, actual)
	}
}

func TestVersionPrerelease(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		version  string
		expected string
	}{
		{"1.2.3", ""},
		{"1.2-beta", "beta"},
		{"1.2.0-x.Y.0", "x.Y.0"},
		{"1.2.0-x.Y.0+metadata", "x.Y.0"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "metadata-1.2.0"},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		assert.Nil(err)

		actual := v.Prerelease()
		expected := tc.expected
		assert.Equal(expected, actual)
	}
}

func TestVersionSegments(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		version  string
		expected []int
	}{
		{"1.2.3", []int{1, 2, 3}},
		{"1.2-beta", []int{1, 2, 0}},
		{"1-x.Y.0", []int{1, 0, 0}},
		{"1.2.0-x.Y.0+metadata", []int{1, 2, 0}},
		{"1.2.0-metadata-1.2.0+metadata~dist", []int{1, 2, 0}},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		assert.Nil(err)

		actual := v.Segments()
		expected := tc.expected
		assert.Equal(expected, actual)
	}
}

func TestVersionSegments64(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		version  string
		expected []int64
	}{
		{"1.2.3", []int64{1, 2, 3}},
		{"1.2-beta", []int64{1, 2, 0}},
		{"1-x.Y.0", []int64{1, 0, 0}},
		{"1.2.0-x.Y.0+metadata", []int64{1, 2, 0}},
		{"1.4.9223372036854775807", []int64{1, 4, 9223372036854775807}},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		assert.Nil(err)

		actual := v.Segments64()
		expected := tc.expected
		assert.Equal(expected, actual)
	}
}

func TestVersionString(t *testing.T) {
	assert := assert.New(t)

	cases := [][]string{
		{"1.2.3", "1.2.3"},
		{"1.2-beta", "1.2.0-beta"},
		{"1.2.0-x.Y.0", "1.2.0-x.Y.0"},
		{"1.2.0-x.Y.0+metadata", "1.2.0-x.Y.0+metadata"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "1.2.0-metadata-1.2.0+metadata~dist"},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc[0])
		assert.Nil(err)

		actual := v.String()
		expected := tc[1]
		assert.Equal(expected, actual)
	}
}

func TestCollection(t *testing.T) {
	assert := assert.New(t)

	versionsRaw := []string{
		"1.1.1",
		"1.0",
		"1.2",
		"2",
		"0.7.1",
	}

	versions := make([]*Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, err := NewVersion(raw)
		assert.Nil(err)
		versions[i] = v
	}

	sort.Sort(Collection(versions))

	actual := make([]string, len(versions))
	for i, v := range versions {
		actual[i] = v.String()
	}

	expected := []string{
		"0.7.1",
		"1.0.0",
		"1.1.1",
		"1.2.0",
		"2.0.0",
	}

	assert.Equal(expected, actual)
}
