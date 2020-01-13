package web

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/blend/go-sdk/bufferutil"
	"github.com/blend/go-sdk/ex"
	templatehelpers "github.com/blend/go-sdk/template"
)

const (
	// DefaultTemplateNameBadRequest is the default template name for bad request view results.
	DefaultTemplateNameBadRequest = "bad_request"
	// DefaultTemplateNameInternalError is the default template name for internal server error view results.
	DefaultTemplateNameInternalError = "error"
	// DefaultTemplateNameNotFound is the default template name for not found error view results.
	DefaultTemplateNameNotFound = "not_found"
	// DefaultTemplateNameNotAuthorized is the default template name for not authorized error view results.
	DefaultTemplateNameNotAuthorized = "not_authorized"
	// DefaultTemplateNameStatus is the default template name for status view results.
	DefaultTemplateNameStatus = "status"

	// DefaultTemplateBadRequest is a basic view.
	DefaultTemplateBadRequest = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Bad Request</h4></body><pre>{{ .ViewModel }}</pre></html>`
	// DefaultTemplateInternalError is a basic view.
	DefaultTemplateInternalError = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Internal Error</h4><pre>{{ .ViewModel }}</body></html>`
	// DefaultTemplateNotAuthorized is a basic view.
	DefaultTemplateNotAuthorized = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Not Authorized</h4></body></html>`
	// DefaultTemplateNotFound is a basic view.
	DefaultTemplateNotFound = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Not Found</h4></body></html>`
	// DefaultTemplateStatus is a basic view.
	DefaultTemplateStatus = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>{{ .ViewModel.StatusCode }}</h4></body><pre>{{ .ViewModel.Response }}</pre></html>`
)

// Assert the view cache is a result provider.
var (
	_ ResultProvider = (*ViewCache)(nil)
)

// NewViewCache returns a new view cache.
func NewViewCache(options ...ViewCacheOption) *ViewCache {
	vc := &ViewCache{
		FuncMap:                   template.FuncMap(templatehelpers.ViewFuncs{}.FuncMap()),
		BufferPool:                bufferutil.NewPool(1024),
		InternalErrorTemplateName: DefaultTemplateNameInternalError,
		BadRequestTemplateName:    DefaultTemplateNameBadRequest,
		NotFoundTemplateName:      DefaultTemplateNameNotFound,
		NotAuthorizedTemplateName: DefaultTemplateNameNotAuthorized,
		StatusTemplateName:        DefaultTemplateNameStatus,
	}
	for _, option := range options {
		option(vc)
	}
	return vc
}

// ViewCache is the cached views used in view results.
type ViewCache struct {
	sync.Mutex
	LiveReload bool
	FuncMap    template.FuncMap
	Paths      []string
	Literals   []string
	Templates  *template.Template
	BufferPool *bufferutil.Pool

	BadRequestTemplateName    string
	InternalErrorTemplateName string
	NotFoundTemplateName      string
	NotAuthorizedTemplateName string
	StatusTemplateName        string
}

// Initialize caches templates by path.
func (vc *ViewCache) Initialize() error {
	vc.Lock()
	defer vc.Unlock()
	if vc.Templates == nil && !vc.LiveReload {
		return vc.initialize()
	}
	return nil
}

// Parse parses the view tree.
func (vc *ViewCache) Parse() (views *template.Template, err error) {
	views = template.New("").Funcs(vc.FuncMap)
	if len(vc.Paths) > 0 {
		views, err = views.ParseFiles(vc.Paths...)
		if err != nil {
			err = ex.New(err)
			return
		}
	}

	if len(vc.Literals) > 0 {
		for _, viewLiteral := range vc.Literals {
			views, err = views.Parse(viewLiteral)
			if err != nil {
				err = ex.New(err)
				return
			}
		}
	}
	return
}

// Lookup looks up a view.
func (vc *ViewCache) Lookup(name string) (*template.Template, error) {
	if vc.Templates == nil {
		templates, err := vc.Parse()
		if err != nil {
			return nil, err
		}
		return templates.Lookup(name), nil
	}
	return vc.Templates.Lookup(name), nil
}

// ----------------------------------------------------------------------
// results
// ----------------------------------------------------------------------

// BadRequest returns a view result.
func (vc *ViewCache) BadRequest(err error) Result {
	t, viewErr := vc.Lookup(vc.BadRequestTemplateName)
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("default").Parse(DefaultTemplateBadRequest)
	}

	return &ViewResult{
		ViewName:   vc.BadRequestTemplateName,
		StatusCode: http.StatusBadRequest,
		ViewModel:  err,
		Template:   t,
		Views:      vc,
	}
}

// InternalError returns a view result.
func (vc *ViewCache) InternalError(err error) Result {
	t, viewErr := vc.Lookup(vc.InternalErrorTemplateName)
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateInternalError)
	}
	return ResultWithLoggedError(&ViewResult{
		ViewName:   vc.InternalErrorTemplateName,
		StatusCode: http.StatusInternalServerError,
		ViewModel:  err,
		Template:   t,
		Views:      vc,
	}, err)
}

// NotFound returns a view result.
func (vc *ViewCache) NotFound() Result {
	t, viewErr := vc.Lookup(vc.NotFoundTemplateName)
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateNotFound)
	}
	return &ViewResult{
		ViewName:   vc.NotFoundTemplateName,
		StatusCode: http.StatusNotFound,
		Template:   t,
		Views:      vc,
	}
}

// NotAuthorized returns a view result.
func (vc *ViewCache) NotAuthorized() Result {
	t, err := vc.Lookup(vc.NotAuthorizedTemplateName)
	if err != nil {
		return vc.viewError(err)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateNotAuthorized)
	}

	return &ViewResult{
		ViewName:   vc.NotAuthorizedTemplateName,
		StatusCode: http.StatusUnauthorized,
		Template:   t,
		Views:      vc,
	}
}

// Status returns a status view result.
func (vc *ViewCache) Status(statusCode int, response ...interface{}) Result {
	t, viewErr := vc.Lookup(vc.StatusTemplateName)
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateStatus)
	}

	return &ViewResult{
		Views:      vc,
		ViewName:   vc.StatusTemplateName,
		StatusCode: statusCode,
		Template:   t,
		ViewModel:  StatusViewModel{StatusCode: statusCode, Response: ResultOrDefault(http.StatusText(statusCode), response...)},
	}
}

// View returns a view result.
func (vc *ViewCache) View(viewName string, viewModel interface{}) Result {
	return vc.ViewStatus(http.StatusOK, viewName, viewModel)
}

// ViewStatus returns a view result with a given status code..
func (vc *ViewCache) ViewStatus(statusCode int, viewName string, viewModel interface{}) Result {
	t, err := vc.Lookup(viewName)
	if err != nil {
		return vc.viewError(err)
	}
	if t == nil {
		return vc.InternalError(ex.New(ErrUnsetViewTemplate, ex.OptMessagef("viewname: %s", viewName)))
	}

	return &ViewResult{
		ViewName:   viewName,
		StatusCode: statusCode,
		ViewModel:  viewModel,
		Template:   t,
		Views:      vc,
	}
}

// ----------------------------------------------------------------------
// properties
// ----------------------------------------------------------------------

// AddPaths adds paths to the view collection.
func (vc *ViewCache) AddPaths(paths ...string) {
	vc.Paths = append(vc.Paths, paths...)
}

// AddLiterals adds view literal strings to the view collection.
func (vc *ViewCache) AddLiterals(views ...string) {
	vc.Literals = append(vc.Literals, views...)
}

// ----------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------

func (vc *ViewCache) viewError(err error) Result {
	t, _ := template.New("").Parse(DefaultTemplateInternalError)
	return &ViewResult{
		ViewName:   DefaultTemplateNameInternalError,
		StatusCode: http.StatusInternalServerError,
		ViewModel:  err,
		Template:   t,
		Views:      vc,
	}
}

func (vc *ViewCache) initialize() error {
	if len(vc.Paths) == 0 && len(vc.Literals) == 0 {
		return nil
	}
	views, err := vc.Parse()
	if err != nil {
		return err
	}
	vc.Templates = views
	return nil
}
