/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// Float64Ptr returns an Float64Source for a given float64 pointer.
func Float64Ptr(value *float64) Float64Source {
	return Float64PtrSource{Value: value}
}

var (
	_ Float64Source = (*Float64PtrSource)(nil)
)

// Float64PtrSource is a Float64Source that wraps a float64 pointer.
type Float64PtrSource struct {
	Value *float64
}

// Float64 implements Float64Source.
func (fps Float64PtrSource) Float64(_ context.Context) (*float64, error) {
	return fps.Value, nil
}
