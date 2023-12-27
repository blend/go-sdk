//go:build darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd
// +build darwin dragonfly freebsd illumos linux netbsd openbsd

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package filelock

import (
	"io/fs"
	"syscall"
)

type lockType int16

const (
	readLock	lockType	= syscall.LOCK_SH
	writeLock	lockType	= syscall.LOCK_EX
)

func lock(f File, lt lockType) (err error) {
	for {
		err = syscall.Flock(int(f.Fd()), int(lt))
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		return &fs.PathError{
			Op:	lt.String(),
			Path:	f.Name(),
			Err:	err,
		}
	}
	return nil
}

func unlock(f File) error {
	return lock(f, syscall.LOCK_UN)
}

func isNotSupported(err error) bool {
	return err == syscall.ENOSYS || err == syscall.ENOTSUP || err == syscall.EOPNOTSUPP || err == ErrNotSupported
}
