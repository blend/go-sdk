/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// Int32Source is a type that can return a value.
type Int32Source interface {
	// Int32 should return an int32 if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Int32(ctx context.Context) (*int32, error)
}
