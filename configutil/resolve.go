/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
)

// ResolveAction is a step in resolution.
type ResolveAction func(context.Context) error

// Resolve returns the first non-nil error in a list.
func Resolve(ctx context.Context, steps ...ResolveAction) (err error) {
	for _, step := range steps {
		if err = step(ctx); err != nil {
			return err
		}
	}
	return nil
}
