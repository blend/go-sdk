/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cache

import "sync"

// Cache is a type that implements the cache interface.
type Cache interface {
	Has(key interface{}) bool
	GetOrSet(key interface{}, valueProvider func() (interface{}, error), options ...ValueOption) (interface{}, bool, error)
	Set(key, value interface{}, options ...ValueOption)
	Get(key interface{}) (interface{}, bool)
	Remove(key interface{}) (interface{}, bool)
}

// Locker is a cache type that supports external control of locking for both exclusive and reader/writer locks.
type Locker interface {
	sync.Locker
	RLock()
	RUnlock()
}
