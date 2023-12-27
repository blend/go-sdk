/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import "sync"

var (
	_default	*JobManager
	_defaultLock	sync.Mutex
)

// Default returns a shared instance of a JobManager.
// If unset, it will initialize it with `New()`.
func Default() *JobManager {
	if _default == nil {
		_defaultLock.Lock()
		defer _defaultLock.Unlock()

		if _default == nil {
			_default = New()
		}
	}
	return _default
}

// SetDefault sets the default job manager.
func SetDefault(jm *JobManager) {
	_defaultLock.Lock()
	_default = jm
	_defaultLock.Unlock()
}
