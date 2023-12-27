/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

// Int64Ptr returns an Int64Source for a given Int64 poInt64er.
func Int64Ptr(value *int64) Int64Source {
	return Int64PtrSource{Value: value}
}

var (
	_ Int64Source = (*Int64PtrSource)(nil)
)

// Int64PtrSource is a Int64Source that wraps an Int64 poInt64er.
type Int64PtrSource struct {
	Value *int64
}

// Int64 implements Int64Source.
func (ips Int64PtrSource) Int64(_ context.Context) (*int64, error) {
	return ips.Value, nil
}
