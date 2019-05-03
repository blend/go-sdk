package template

import (
	"bytes"
	"io"
	"io/ioutil"

	texttemplate "text/template"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
)

// New creates a new template.
func New() *Template {
	temp := &Template{
		Viewmodel: Viewmodel{
			vars: Vars{},
			env:  env.Env(),
		},
	}
	temp.funcs = texttemplate.FuncMap(ViewFuncs{}.FuncMap())
	return temp
}

// NewFromFile creates a new template from a file.
func NewFromFile(filepath string) (*Template, error) {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, ex.New(err)
	}

	return New().WithName(filepath).WithBody(string(contents)), nil
}

// Template is a wrapper for html.Template.
type Template struct {
	Viewmodel
	name       string
	body       string
	includes   []string
	funcs      texttemplate.FuncMap
	leftDelim  string
	rightDelim string
}

// WithName sets the template name.
func (t *Template) WithName(name string) *Template {
	t.name = name
	return t
}

// Name returns the template name if set, or if not set, just "template" as a constant.
func (t *Template) Name() string {
	if len(t.name) > 0 {
		return t.name
	}
	return "template"
}

// WithDelims sets the template action delimiters, treating empty string as default delimiter.
func (t *Template) WithDelims(left, right string) *Template {
	t.leftDelim = left
	t.rightDelim = right
	return t
}

// WithBody sets the template body and returns a reference to the template object.
func (t *Template) WithBody(body string) *Template {
	t.body = body
	return t
}

// WithInclude includes a (sub) template into the rendering assets.
func (t *Template) WithInclude(body string) *Template {
	t.includes = append(t.includes, body)
	return t
}

// Body returns the template body.
func (t *Template) Body() string {
	return t.body
}

// WithVar sets a variable and returns a reference to the template object.
func (t *Template) WithVar(key string, value interface{}) *Template {
	t.SetVar(key, value)
	return t
}

// SetVar sets a var in the template.
func (t *Template) SetVar(key string, value interface{}) {
	t.vars[key] = value
}

// WithVars reads a map of variables into the template.
func (t *Template) WithVars(vars Vars) *Template {
	t.vars = MergeVars(t.vars, vars)
	return t
}

// WithEnvVars sets the environment variables.
func (t *Template) WithEnvVars(envVars env.Vars) *Template {
	t.Viewmodel.env = env.Merge(t.Viewmodel.env, envVars)
	return t
}

// SetVarsFromFile reads vars from a file and merges them
// with the current variables set.
func (t *Template) SetVarsFromFile(path string) error {
	fileVars, err := NewVarsFromPath(path)
	if err != nil {
		return err
	}

	t.vars = MergeVars(t.vars, fileVars)
	return nil
}

// Process processes the template.
func (t *Template) Process(dst io.Writer) error {
	base := texttemplate.New(t.Name()).Funcs(t.ViewFuncs()).Delims(t.leftDelim, t.rightDelim)

	var err error
	for _, include := range t.includes {
		_, err = base.New(t.Name()).Parse(include)
		if err != nil {
			return err
		}
	}

	final, err := base.New(t.Name()).Parse(t.body)
	if err != nil {
		return err
	}
	return final.Execute(dst, t.Viewmodel)
}

// ProcessString is a helper to process the template as a string.
func (t *Template) ProcessString() (string, error) {
	buffer := new(bytes.Buffer)
	err := t.Process(buffer)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// MustProcessString is a helper to process a template as a string
// and panic on error.
func (t *Template) MustProcessString() string {
	output, err := t.ProcessString()
	if err != nil {
		panic(err)
	}
	return output
}

// ViewFuncs returns the view funcs.
func (t *Template) ViewFuncs() texttemplate.FuncMap {
	return t.funcs
}
