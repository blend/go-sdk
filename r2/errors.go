package r2

import (
	"github.com/blend/go-sdk/ex"
)

// Error Constants
const (
	ErrNoContentJSON ex.Class = "server returned an http 204 for a request expecting json"
	ErrNoContentXML  ex.Class = "server returned an http 204 for a request expecting xml"
)
