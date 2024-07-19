/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import "github.com/blend/go-sdk/graceful"

// Validate the interface is satisfied.
var (
	_ (graceful.Graceful) = (*Graceful)(nil)
)
