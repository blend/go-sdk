/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

var (
	_ BoolSource = (*BoolFunc)(nil)
)

// BoolFunc is a bool value source.
// It can be used with configutil.SetBool
type BoolFunc func(context.Context) (*bool, error)

// Bool returns an invocation of the function.
func (vf BoolFunc) Bool(ctx context.Context) (*bool, error) {
	return vf(ctx)
}
