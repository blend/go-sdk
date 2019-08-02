package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

// NewStaticFileServer returns a new static file cache.
func NewStaticFileServer(searchPaths ...http.FileSystem) *StaticFileServer {
	return &StaticFileServer{
		SearchPaths: searchPaths,
	}
}

// StaticFileServer is a cache of static files.
type StaticFileServer struct {
	sync.RWMutex
	SearchPaths   []http.FileSystem
	RewriteRules  []RewriteRule
	Middleware    []Middleware
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
// It will run middleware if specified before serving the file.
func (sc *StaticFileServer) Action(r *Ctx) Result {
	filePath, err := r.RouteParam("filepath")
	if err != nil {
		return r.DefaultProvider.BadRequest(err)
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

// ServeFile writes the file to the response without running middleware.
func (sc *StaticFileServer) ServeFile(r *Ctx, filePath string) Result {
	f, err := sc.ResolveFile(filePath)
	if f == nil || (err != nil && os.IsNotExist(err)) {
		return r.DefaultProvider.NotFound()
	}
	if err != nil {
		return r.DefaultProvider.InternalError(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return r.DefaultProvider.InternalError(err)
	}
	http.ServeContent(r.Response, r.Request, filePath, finfo.ModTime(), f)
	return nil
}

// ServeCachedFile writes the file to the response.
func (sc *StaticFileServer) ServeCachedFile(r *Ctx, filepath string) Result {
	file, err := sc.ResolveCachedFile(filepath)
	if err != nil {
		return r.DefaultProvider.InternalError(err)
	}
	if file == nil {
		http.NotFound(r.Response, r.Request)
		return nil
	}
	http.ServeContent(r.Response, r.Request, filepath, file.ModTime, file.Contents)
	return nil
}

// ResolveFile resolves a file from rewrite rules and search paths.
func (sc *StaticFileServer) ResolveFile(filePath string) (f http.File, err error) {
	for _, rule := range sc.RewriteRules {
		if matched, newFilePath := rule.Apply(filePath); matched {
			filePath = newFilePath
		}
	}

	// for each searchpath, sniff if the file exists ...
	var openErr error
	for _, searchPath := range sc.SearchPaths {
		f, openErr = searchPath.Open(filePath)
		if openErr == nil {
			break
		}
	}
	if openErr != nil && !os.IsNotExist(openErr) {
		err = openErr
		return
	}
	return
}

// ResolveCachedFile returns a cached file at a given path.
// It returns the cached instance of a file if it exists, and adds it to the cache if there is a miss.
func (sc *StaticFileServer) ResolveCachedFile(filepath string) (*CachedStaticFile, error) {
	sc.RLock()
	if sc.Cache != nil {
		if file, ok := sc.Cache[filepath]; ok {
			sc.RUnlock()
			return file, nil
		}
	}
	sc.RUnlock()

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
		Size:     len(contents),
	}

	sc.Cache[filepath] = file
	return file, nil
}
