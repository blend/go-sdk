/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package semver

import "github.com/blend/go-sdk/ex"

const (
	// ErrConstraintFailed is returned by validators.
	ErrConstraintFailed ex.Class = "semver; constraint failed"
)
