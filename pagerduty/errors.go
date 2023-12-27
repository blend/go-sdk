/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import "github.com/blend/go-sdk/ex"

// Errors
const (
	ErrNon200Status	ex.Class	= "non-200 status code from remote"
	Err404Status	ex.Class	= "404 status code from remote"
)
