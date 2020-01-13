package bindata

import (
	"compress/gzip"
	"crypto/md5"
	"hash"
	"io"

	"github.com/blend/go-sdk/ex"
)

// NewFileCompressor returns a new file compressor.
func NewFileCompressor(src io.ReadCloser) *FileCompressor {
	return &FileCompressor{
		Source: src,
		MD5:    md5.New(),
	}
}

// FileCompressor reads a file an returns compressed output.
type FileCompressor struct {
	Source io.ReadCloser
	MD5    hash.Hash
}

// WriteTo copies the source to the destination as a compressed output.
// It also sums the source into the MD5 hash (uncompressed)>
func (fc *FileCompressor) WriteTo(dst io.Writer) (written int64, err error) {
	gzw := gzip.NewWriter(dst)
	defer gzw.Close()

	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := fc.Source.Read(buf)
		if nr > 0 {
			_, em := fc.MD5.Write(buf[0:nr])
			if em != nil {
				err = ex.New(em)
			}
			nw, ew := gzw.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ex.New(ew)
				break
			}
			if nr != nw {
				err = ex.New(io.ErrShortWrite)
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = ex.New(er)
			}
			break
		}
	}
	return written, err
}

// Close closes the source stream.
func (fc *FileCompressor) Close() error {
	return ex.New(fc.Source.Close())
}
