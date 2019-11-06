package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/blend/go-sdk/webutil"
)

// Interface assertions
var (
	_ Fileserver = (*CachedStaticFileServer)(nil)
)

// NewCachedStaticFileServer returns a new static file cache.
func NewCachedStaticFileServer(searchPaths ...http.FileSystem) *CachedStaticFileServer {
	return &CachedStaticFileServer{
		StaticFileServer: NewStaticFileServer(searchPaths...),
		files:            map[string]*CachedStaticFile{},
	}
}

// CachedStaticFileServer  is a cache of static files.
type CachedStaticFileServer struct {
	sync.Mutex
	*StaticFileServer

	files map[string]*CachedStaticFile
}

// Files returns the underlying file cache.
// Pragma; this should only be used in debugging, as during runtime locks are required to interact with this cache.
func (csfs *CachedStaticFileServer) Files() map[string]*CachedStaticFile {
	return csfs.files
}

// Action is the entrypoint for the static server.
func (csfs *CachedStaticFileServer) Action(r *Ctx) Result {
	filepath, err := r.RouteParam("filepath")
	if err != nil {
		return r.DefaultResultProvider().BadRequest(err)
	}

	if csfs.middleware != nil {
		return csfs.middleware(r)
	}
	return csfs.ServeFile(r, filepath)
}

// File returns a cached file at a given path.
// It returns the cached instance of a file if it exists, and adds it to the cache if there is a miss.
func (csfs *CachedStaticFileServer) File(filepath string) (*CachedStaticFile, error) {
	csfs.Lock()
	defer csfs.Unlock()

	if file, ok := csfs.files[filepath]; ok {
		return file, nil
	}

	diskFile, err := csfs.StaticFileServer.ResolveFile(filepath)
	if err != nil {
		return nil, err
	}

	finfo, err := diskFile.Stat()
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(diskFile)
	if err != nil {
		return nil, err
	}

	file := &CachedStaticFile{
		Path:     filepath,
		ETag:     webutil.ETag(contents),
		Contents: bytes.NewReader(contents),
		ModTime:  finfo.ModTime(),
		Size:     len(contents),
	}

	csfs.files[filepath] = file
	return file, nil
}

// ServeFile writes the file to the response.
func (csfs *CachedStaticFileServer) ServeFile(r *Ctx, filepath string) Result {
	for key, values := range csfs.StaticFileServer.headers {
		for _, value := range values {
			r.Response().Header().Set(key, value)
		}
	}
	file, err := csfs.File(filepath)
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}
	_ = file.Render(r)
	return nil
}
