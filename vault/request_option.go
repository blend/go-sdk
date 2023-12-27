/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import (
	"strconv"

	"github.com/blend/go-sdk/webutil"
)

// CallOption a thing that we can do to modify a request.
type CallOption = webutil.RequestOption

// OptVersion adds a version to the request.
func OptVersion(version int) CallOption {
	return webutil.OptQueryValue("version", strconv.Itoa(version))
}

// OptList adds a list parameter to the request.
func OptList() CallOption {
	return webutil.OptQueryValue("list", "true")
}
