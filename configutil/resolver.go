/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// Resolver is a type that can be resolved.
type Resolver interface {
	Resolve(context.Context) error
}
