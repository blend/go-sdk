package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/blend/go-sdk/webutil"
)

// NewStaticFileServer returns a new static file cache.
func NewStaticFileServer(options ...StaticFileserverOption) *StaticFileServer {
	var sfs StaticFileServer
	for _, opt := range options {
		opt(&sfs)
	}
	return &sfs
}

// StaticFileserverOption are options for static fileservers.
type StaticFileserverOption func(*StaticFileServer)

// OptStaticFileServerSearchPaths sets the static fileserver search paths.
func OptStaticFileServerSearchPaths(searchPaths ...http.FileSystem) StaticFileserverOption {
	return func(sfs *StaticFileServer) {
		sfs.SearchPaths = searchPaths
	}
}

// OptStaticFileServerHeaders sets the static fileserver default headers..
func OptStaticFileServerHeaders(headers http.Header) StaticFileserverOption {
	return func(sfs *StaticFileServer) {
		sfs.Headers = headers
	}
}

// OptStaticFileServerCacheDisabled sets the static fileserver should read from disk for each request.
func OptStaticFileServerCacheDisabled(cacheDisabled bool) StaticFileserverOption {
	return func(sfs *StaticFileServer) {
		sfs.CacheDisabled = cacheDisabled
	}
}

// StaticFileServer is a cache of static files.
// It can operate in cached mode, or with `CacheDisabled` set to `true`
// it will read from disk for each request.
// In cached mode, it automatically adds etags for files it caches.
type StaticFileServer struct {
	sync.RWMutex

	SearchPaths   []http.FileSystem
	RewriteRules  []RewriteRule
	Headers       http.Header
	CacheDisabled bool
	Cache         map[string]*CachedStaticFile
}

// AddHeader adds a header to the static cache results.
func (sc *StaticFileServer) AddHeader(key, value string) {
	if sc.Headers == nil {
		sc.Headers = http.Header{}
	}
	sc.Headers[key] = append(sc.Headers[key], value)
}

// AddRewriteRule adds a static re-write rule.
// This is meant to modify the path of a file from what is requested by the browser
// to how a file may actually be accessed on disk.
// Typically re-write rules are used to enforce caching semantics.
func (sc *StaticFileServer) AddRewriteRule(match string, action RewriteAction) error {
	expr, err := regexp.Compile(match)
	if err != nil {
		return err
	}
	sc.RewriteRules = append(sc.RewriteRules, RewriteRule{
		MatchExpression: match,
		expr:            expr,
		Action:          action,
	})
	return nil
}

// Action is the entrypoint for the static server.
// It  adds default headers if specified, and then serves the file from disk
// or from a pull-through cache if enabled.
func (sc *StaticFileServer) Action(r *Ctx) Result {
	filePath, err := r.RouteParam("filepath")
	if err != nil {
		if r.DefaultProvider != nil {
			return r.DefaultProvider.BadRequest(err)
		}
		http.Error(r.Response, err.Error(), http.StatusBadRequest)
		return nil
	}

	for key, values := range sc.Headers {
		for _, value := range values {
			r.Response.Header().Set(key, value)
		}
	}

	if sc.CacheDisabled {
		return sc.ServeFile(r, filePath)
	}
	return sc.ServeCachedFile(r, filePath)
}

// ServeFile writes the file to the response by reading from disk
// for each request (i.e. skipping the cache)
func (sc *StaticFileServer) ServeFile(r *Ctx, filePath string) Result {
	f, err := sc.ResolveFile(filePath)
	if err != nil {
		return sc.fileError(r, err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return sc.fileError(r, err)
	}
	http.ServeContent(r.Response, r.Request, filePath, finfo.ModTime(), f)
	return nil
}

// ServeCachedFile writes the file to the response, potentially
// serving a cached instance of the file.
func (sc *StaticFileServer) ServeCachedFile(r *Ctx, filepath string) Result {
	file, err := sc.ResolveCachedFile(filepath)
	if err != nil {
		return sc.fileError(r, err)
	}
	if file.ETag != "" {
		r.Response.Header().Set(webutil.HeaderETag, file.ETag)
	}
	http.ServeContent(r.Response, r.Request, filepath, file.ModTime, file.Contents)
	return nil
}

// ResolveFile resolves a file from rewrite rules and search paths.
// First the file path is modified according to the rewrite rules.
// Then each search path is checked for the resolved file path.
func (sc *StaticFileServer) ResolveFile(filePath string) (f http.File, err error) {
	for _, rule := range sc.RewriteRules {
		if matched, newFilePath := rule.Apply(filePath); matched {
			filePath = newFilePath
		}
	}
	for _, searchPath := range sc.SearchPaths {
		f, err = searchPath.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return
		}
	}
	return
}

// ResolveCachedFile returns a cached file at a given path.
// It returns the cached instance of a file if it exists, and adds it to the cache if there is a miss.
func (sc *StaticFileServer) ResolveCachedFile(filepath string) (*CachedStaticFile, error) {
	// start in read shared mode
	sc.RLock()
	if sc.Cache != nil {
		if file, ok := sc.Cache[filepath]; ok {
			sc.RUnlock()
			return file, nil
		}
	}
	sc.RUnlock()

	// transition to exclusive write mode
	sc.Lock()
	defer sc.Unlock()

	if sc.Cache == nil {
		sc.Cache = make(map[string]*CachedStaticFile)
	}
	// double check ftw
	if file, ok := sc.Cache[filepath]; ok {
		return file, nil
	}

	diskFile, err := sc.ResolveFile(filepath)
	if err != nil {
		return nil, err
	}

	if diskFile == nil {
		sc.Cache[filepath] = nil
		return nil, nil
	}

	finfo, err := diskFile.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	contents, err := ioutil.ReadAll(diskFile)
	if err != nil {
		return nil, err
	}

	file := &CachedStaticFile{
		Path:     filepath,
		Contents: bytes.NewReader(contents),
		ModTime:  finfo.ModTime(),
		ETag:     webutil.ETag(contents),
		Size:     len(contents),
	}

	sc.Cache[filepath] = file
	return file, nil
}

func (sc *StaticFileServer) fileError(r *Ctx, err error) Result {
	if os.IsNotExist(err) {
		if r.DefaultProvider != nil {
			return r.DefaultProvider.NotFound()
		}
		http.NotFound(r.Response, r.Request)
		return nil
	}
	if r.DefaultProvider != nil {
		return r.DefaultProvider.InternalError(err)
	}
	http.Error(r.Response, err.Error(), http.StatusInternalServerError)
	return nil
}
