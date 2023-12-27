/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

import (
	"os"
	"sync"
)

var (
	_env     Vars
	_envLock = sync.Mutex{}
)

// Env returns the current env var set.
func Env() Vars {
	if _env == nil {
		_envLock.Lock()
		defer _envLock.Unlock()
		if _env == nil {
			_env = New(OptEnviron(os.Environ()...))
		}
	}
	return _env
}

// SetEnv sets the env vars.
func SetEnv(vars Vars) {
	_envLock.Lock()
	_env = vars
	_envLock.Unlock()
}

// Restore sets .Env() to the current os environment.
func Restore() {
	SetEnv(New(OptEnviron(os.Environ()...)))
}

// Clear sets .Env() to an empty env var set.
func Clear() {
	SetEnv(New())
}
