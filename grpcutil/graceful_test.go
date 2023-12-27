/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import "github.com/blend/go-sdk/graceful"

// Validate the interface is satisfied.
var (
	_ (graceful.Graceful) = (*Graceful)(nil)
)
