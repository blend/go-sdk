/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/stringutil"
)

// CopyAll copies all files and directories from a source path to a destination path recurrsively.
func CopyAll(destination, source string, opts ...CopyAllOption) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	var finalOptions CopyAllOptions
	for _, opt := range opts {
		opt(&finalOptions)
	}
	return copyAll(destination, source, info, finalOptions)
}

// OptCopyAllSymlinkMode sets the symlink mode.
func OptCopyAllSymlinkMode(mode CopyAllSymlinkMode) CopyAllOption {
	return func(cop *CopyAllOptions) { cop.SymlinkMode = mode }
}

// OptCopyAllSkipGlobs sets the skip provider to a glob matcher based on a given set of glob(s).
func OptCopyAllSkipGlobs(globs ...string) CopyAllOption {
	return func(cop *CopyAllOptions) {
		cop.SkipProvider = func(fileInfo os.FileInfo) bool {
			for _, glob := range globs {
				if stringutil.Glob(fileInfo.Name(), glob) {
					return true
				}
			}
			return false
		}
	}
}

// CopyAllOptions are the options for copy all.
type CopyAllOptions struct {
	SymlinkMode	CopyAllSymlinkMode
	SkipProvider	func(os.FileInfo) bool
}

// CopyAllOption is a mutator for copy all options
type CopyAllOption func(*CopyAllOptions)

// CopyAllSymlinkMode is how symlinks should be handled
type CopyAllSymlinkMode int

// CopyAllSymlinkMode(s)
var (
	// CopyAllSymlinkModeShallow will copy links from the source to the destination as links.
	CopyAllSymlinkModeShallow	CopyAllSymlinkMode	= 0
	// CopyAllSymlinkModeDeep will traverse into the link destination and copy any files recursively.
	CopyAllSymlinkModeDeep	CopyAllSymlinkMode	= 1
	// CopyAllSymlinkModeSkip will skip any links discovered.
	CopyAllSymlinkModeSkip	CopyAllSymlinkMode	= 2
)

// copyAll switches proper copy functions regarding file type, etc...
// If there would be anything else here, add a case to this switchboard.
func copyAll(destination, source string, info os.FileInfo, opts CopyAllOptions) error {
	if opts.SkipProvider != nil {
		if opts.SkipProvider(info) {
			return nil
		}
	}

	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return symCopy(destination, source, info, opts)
	case info.IsDir():
		return dirCopy(destination, source, info, opts)
	default:
		return fileCopy(destination, source, info, opts)
	}
}

func symCopy(destination, source string, sourceDirInfo os.FileInfo, opts CopyAllOptions) error {
	switch opts.SymlinkMode {
	case CopyAllSymlinkModeShallow:
		return linkCopy(destination, source)
	case CopyAllSymlinkModeDeep:
		orig, err := os.Readlink(source)
		if err != nil {
			return err
		}
		originalSourceDirInfo, err := os.Lstat(orig)
		if err != nil {
			return err
		}
		return copyAll(destination, orig, originalSourceDirInfo, opts)
	case CopyAllSymlinkModeSkip:
		fallthrough
	default:
		return nil	// do nothing
	}
}

// linkCopy copies a symlink.
func linkCopy(destination, source string) error {
	src, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(src, destination)
}

// dirCopy copies a directory recursively
func dirCopy(destination, source string, sourceDirInfo os.FileInfo, opts CopyAllOptions) (err error) {
	if _, statErr := os.Stat(destination); statErr != nil {
		// create the destination with the source permissions
		if err = os.MkdirAll(destination, sourceDirInfo.Mode()); err != nil {
			return
		}
	}

	var sourceDirContents []os.DirEntry
	sourceDirContents, err = os.ReadDir(source)
	if err != nil {
		return
	}
	for _, sourceDirItem := range sourceDirContents {
		destinationName := filepath.Join(destination, sourceDirItem.Name())
		sourceName := filepath.Join(source, sourceDirItem.Name())
		var sourceDirItemInfo os.FileInfo
		sourceDirItemInfo, err = sourceDirItem.Info()
		if err != nil {
			return
		}

		if _, statErr := os.Stat(destinationName); statErr != nil {
			if err = copyAll(destinationName, sourceName, sourceDirItemInfo, opts); err != nil {
				return
			}
		}
	}
	return
}

// fileCopy copies a single file
func fileCopy(destination, source string, sourceInfo os.FileInfo, opts CopyAllOptions) (err error) {
	destinationDir := filepath.Dir(destination)
	if _, statErr := os.Stat(destinationDir); statErr != nil {
		if err = os.MkdirAll(destinationDir, sourceInfo.Mode()); err != nil {
			return
		}
	}

	f, err := os.Create(destination)
	if err != nil {
		return
	}
	defer fclose(f, &err)

	s, err := os.Open(source)
	if err != nil {
		return
	}
	defer fclose(s, &err)

	if _, err = io.Copy(f, s); err != nil {
		return
	}
	return
}

// fclose ANYHOW closes file,
// with asiging error raised during Close,
// BUT respecting the error already reported.
func fclose(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}
