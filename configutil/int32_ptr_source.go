/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

// Int32Ptr returns an Int32Source for a given int32 pointer.
func Int32Ptr(value *int32) Int32Source {
	return Int32PtrSource{Value: value}
}

var (
	_ Int32Source = (*Int32PtrSource)(nil)
)

// Int32PtrSource is a Int32Source that wraps an int32 pointer.
type Int32PtrSource struct {
	Value *int32
}

// Int32 implements Int32Source.
func (ips Int32PtrSource) Int32(_ context.Context) (*int32, error) {
	return ips.Value, nil
}
