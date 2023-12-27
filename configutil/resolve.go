/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
