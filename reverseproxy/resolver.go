/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reverseproxy

import (
	"net/http"
	"sync/atomic"
)

// Resolver is a function that takes a request and produces a destination `url.URL`.
type Resolver func(*http.Request, []*Upstream) (*Upstream, error)

// RoundRobinResolver returns a closure based resolver that rotates through upstreams uniformly.
func RoundRobinResolver(upstreams []*Upstream) Resolver {
	if len(upstreams) == 0 {
		return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
			return nil, nil
		}
	}

	if len(upstreams) == 1 {
		return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
			return upstreams[0], nil
		}
	}

	return manyRoundRobinResolver(upstreams)
}

func manyRoundRobinResolver(upstreams []*Upstream) Resolver {
	var index int32
	total := int32(len(upstreams))
	return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
		newIndex := int(atomic.AddInt32(&index, 1) % total)
		return upstreams[newIndex], nil
	}
}
