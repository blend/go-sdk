/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cache

import (
	"time"
)

// ValueOption is an option for a cache value.
type ValueOption func(*Value)

// OptValueTTL sets the ttl for the value.
func OptValueTTL(d time.Duration) ValueOption {
	return func(v *Value) {
		v.Expires = v.Timestamp.Add(d)
	}
}

// OptValueExpires sets the ttl for the value.
func OptValueExpires(expires time.Time) ValueOption {
	return func(v *Value) {
		v.Expires = expires
	}
}

// OptValueTimestamp sets the timestamp for the value.
func OptValueTimestamp(t time.Time) ValueOption {
	return func(v *Value) {
		v.Timestamp = t
	}
}

// OptValueOnRemove sets the on remove handler.
func OptValueOnRemove(handler func(interface{}, RemovalReason)) ValueOption {
	return func(v *Value) {
		v.OnRemove = handler
	}
}

// Value is a cached item.
type Value struct {
	Timestamp	time.Time
	Expires		time.Time
	Key		interface{}
	Value		interface{}
	OnRemove	func(interface{}, RemovalReason)
}
