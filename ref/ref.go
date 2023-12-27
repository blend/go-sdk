/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ref

import "time"

// Ref returns a reference.
func Ref[T any](v T) *T {
	return &v
}

// Deref returns a dereferenced value. For certain types, this could still return nil.
func Deref[T comparable](p *T) T {
	if p == nil {
		var v T
		return v
	}

	return *p
}

// Refs returns a slice of references.
func Refs[T any](values ...T) []*T {
	output := make([]*T, len(values))
	for index := range values {
		output[index] = &values[index]
	}
	return output
}

// Derefs returns a slice of dereferenced values. For certain types, this could still return nil.
func Derefs[T comparable](pointers ...*T) []T {
	output := make([]T, len(pointers))
	for index := range pointers {
		output[index] = Deref(pointers[index])
	}
	return output
}

// String returns a reference.
func String(v string) *string {
	return &v
}

// Strings returns a reference.
func Strings(values ...string) []*string {
	output := make([]*string, len(values))
	for index := range values {
		output[index] = &values[index]
	}
	return output
}

// Bool returns a reference.
func Bool(v bool) *bool {
	return &v
}

// Byte returns a reference.
func Byte(v byte) *byte {
	return &v
}

// Rune returns a reference.
func Rune(v rune) *rune {
	return &v
}

// Uint8 returns a reference.
func Uint8(v uint8) *uint8 {
	return &v
}

// Uint16 returns a reference.
func Uint16(v uint16) *uint16 {
	return &v
}

// Uint32 returns a reference.
func Uint32(v uint32) *uint32 {
	return &v
}

// Uint64 returns a reference.
func Uint64(v uint64) *uint64 {
	return &v
}

// Int8 returns a reference.
func Int8(v int8) *int8 {
	return &v
}

// Int16 returns a reference.
func Int16(v int16) *int16 {
	return &v
}

// Int32 returns a reference.
func Int32(v int32) *int32 {
	return &v
}

// Int64 returns a reference.
func Int64(v int64) *int64 {
	return &v
}

// Int returns a reference.
func Int(v int) *int {
	return &v
}

// Float32 returns a reference.
func Float32(v float32) *float32 {
	return &v
}

// Float64 returns a reference.
func Float64(v float64) *float64 {
	return &v
}

// Time returns a reference.
func Time(v time.Time) *time.Time {
	return &v
}

// Duration returns a reference.
func Duration(v time.Duration) *time.Duration {
	return &v
}
