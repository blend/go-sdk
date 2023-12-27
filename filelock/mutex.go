/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package filelock

import (
	"fmt"
	"io/fs"
	"os"
	"sync"
)

// MutexAt returns a file based mutex at a given path.
func MutexAt(path string) *Mutex {
	futex := &Mutex{
		Path: path,
	}
	return futex
}

// Mutex manages filehandle based locks that are scoped to the lifetime
// of a parent process. If you need semi-durable locks that can span for the duration
// of an action, use `fslock.FSLock` in the github-actions project.
//
// It does not implement `sync.Locker` because file based mutexes can fail to Lock().
type Mutex struct {
	// we still have a sync.Mutex on this to guard against race-prone use cases
	// within the same process.
	mu	sync.Mutex

	// Path is the file path to use as the lock.
	Path	string
}

// RLock creates a reader lock and returns an unlock function.
func (mu *Mutex) RLock() (runlock func(), err error) {
	if mu.Path == "" {
		err = fmt.Errorf("mutex; path unset")
		return
	}

	f, err := mu.openFile(os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	mu.mu.Lock()
	return func() {
		mu.mu.Unlock()
		mu.closeFile(f)
	}, nil
}

// Lock attempts to lock the Mutex.
//
// If successful, Lock returns a non-nil unlock function: it is provided as a
// return-value instead of a separate method to remind the caller to check the
// accompanying error. (See https://golang.org/issue/20803.)
func (mu *Mutex) Lock() (unlock func(), err error) {
	if mu.Path == "" {
		err = fmt.Errorf("mutex; path unset")
		return
	}
	// We could use either O_RDWR or O_WRONLY here. If we choose O_RDWR and the
	// file at mu.Path is write-only, the call to OpenFile will fail with a
	// permission error. That's actually what we want: if we add an RLock method
	// in the future, it should call OpenFile with O_RDONLY and will require the
	// files must be readable, so we should not let the caller make any
	// assumptions about Mutex working with write-only files.
	f, err := mu.openFile(os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	mu.mu.Lock()
	return func() {
		mu.mu.Unlock()
		mu.closeFile(f)
	}, nil
}

func (mu *Mutex) openFile(flag int, perm fs.FileMode) (*os.File, error) {
	// On BSD systems, we could add the O_SHLOCK or O_EXLOCK flag to the OpenFile
	// call instead of locking separately, but we have to support separate locking
	// calls for Linux and Windows anyway, so it's simpler to use that approach
	// consistently.
	f, err := os.OpenFile(mu.Path, flag&^os.O_TRUNC, perm)
	if err != nil {
		return nil, err
	}

	switch flag & (os.O_RDONLY | os.O_WRONLY | os.O_RDWR) {
	case os.O_WRONLY, os.O_RDWR:
		err = Lock(f)
	default:
		err = RLock(f)
	}
	if err != nil {
		f.Close()
		return nil, err
	}

	if flag&os.O_TRUNC == os.O_TRUNC {
		if err := f.Truncate(0); err != nil {
			// The documentation for os.O_TRUNC says “if possible, truncate file when
			// opened”, but doesn't define “possible” (golang.org/issue/28699).
			// We'll treat regular files (and symlinks to regular files) as “possible”
			// and ignore errors for the rest.
			if fi, statErr := f.Stat(); statErr != nil || fi.Mode().IsRegular() {
				Unlock(f)
				f.Close()
				return nil, err
			}
		}
	}
	return f, nil
}

func (mu *Mutex) closeFile(f *os.File) error {
	// Since locking syscalls operate on file descriptors, we must unlock the file
	// while the descriptor is still valid — that is, before the file is closed —
	// and avoid unlocking files that are already closed.
	err := Unlock(f)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}
