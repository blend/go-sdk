/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
