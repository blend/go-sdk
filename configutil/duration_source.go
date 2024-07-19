/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"time"
)

// DurationSource is a type that can return a time.Duration value.
type DurationSource interface {
	// Duration should return a time.Duration if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Duration(context.Context) (*time.Duration, error)
}
