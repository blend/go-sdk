/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package fileutil

import (
	"os"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// NewTemp creates a new temp file with given contents.
func NewTemp(contents []byte) (*Temp, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil, ex.New(err)
	}
	if _, err := f.Write(contents); err != nil {
		return nil, ex.New(err)
	}
	return &Temp{
		file: f,
	}, nil
}

// Temp is a file that deletes itself when closed.
// It does not hold a file handle open, so no
// guarantees are made around the file persisting for the lifetime of the object.
type Temp struct {
	sync.Mutex

	file	*os.File
}

// Name returns the fully qualified file path.
func (tf *Temp) Name() string {
	return tf.file.Name()
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func (tf *Temp) Stat() (os.FileInfo, error) {
	tf.Lock()
	defer tf.Unlock()

	return tf.file.Stat()
}

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (tf *Temp) Read(buffer []byte) (int, error) {
	tf.Lock()
	defer tf.Unlock()

	read, err := tf.file.Read(buffer)
	return read, ex.New(err)
}

// ReadAt reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b).
// At end of file, that error is io.EOF.
func (tf *Temp) ReadAt(buffer []byte, off int64) (int, error) {
	tf.Lock()
	defer tf.Unlock()

	read, err := tf.file.ReadAt(buffer, off)
	return read, ex.New(err)
}

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (tf *Temp) Write(contents []byte) (int, error) {
	tf.Lock()
	defer tf.Unlock()

	written, err := tf.file.Write(contents)
	return written, ex.New(err)
}

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
func (tf *Temp) WriteAt(contents []byte, off int64) (int, error) {
	tf.Lock()
	defer tf.Unlock()

	written, err := tf.file.WriteAt(contents, off)
	return written, ex.New(err)
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (tf *Temp) WriteString(contents string) (int, error) {
	tf.Lock()
	defer tf.Unlock()

	written, err := tf.file.WriteString(contents)
	return written, ex.New(err)
}

// Close closes the file reference and deletes the file.
func (tf *Temp) Close() error {
	tf.Lock()
	defer tf.Unlock()

	if err := tf.file.Close(); err != nil {
		return ex.New(err)
	}
	if err := os.Remove(tf.file.Name()); err != nil {
		return ex.New(err)
	}
	return nil
}
