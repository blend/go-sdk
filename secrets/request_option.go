package secrets

import (
	"strconv"
	"github.com/blend/go-sdk/webutil"
)

// RequestOption a thing that we can do to modify a request.
type RequestOption = webutil.RequestOption

// OptRequestVersion adds a version to the request.
func OptRequestVersion(version int) RequestOption {
	return webutil.OptQueryValue("version", strconv.Itoa(version))
}

// OptRequestList adds a list parameter to the request.
func OptRequestList() RequestOption {
	return webutil.OptQueryValue("list", "true")
}
