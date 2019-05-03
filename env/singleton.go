package env

import "sync"

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
			_env = New(OptFromEnv())
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
	SetEnv(New(OptFromEnv()))
}
