/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package semver

import "github.com/blend/go-sdk/ex"

const (
	// ErrConstraintFailed is returned by validators.
	ErrConstraintFailed ex.Class = "semver; constraint failed"
)
