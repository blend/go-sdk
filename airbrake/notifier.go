package airbrake

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/airbrake/gobrake"
	"github.com/blend/go-sdk/diagnostics"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
)

var (
	_ diagnostics.Notifier = (*Notifier)(nil)
)

var (
	defaultContextOnce sync.Once
	defaultContext     map[string]interface{}
)

// MustNew returns a new notifier and panics on error.
func MustNew(cfg Config) *Notifier {
	notifier, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return notifier
}

// New returns a new notifier.
func New(cfg Config) (*Notifier, error) {
	parsedProjectID, err := strconv.ParseInt(cfg.ProjectID, 10, 64)
	if err != nil {
		return nil, ex.New(err)
	}
	// create a new reporter
	client := gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:   parsedProjectID,
		ProjectKey:  cfg.ProjectKey,
		Environment: cfg.Environment,
	})

	// filter airbrakes from `dev`, `ci`, and `test`.
	client.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		if noticeEnv := notice.Context["environment"]; noticeEnv == env.ServiceEnvDev ||
			noticeEnv == env.ServiceEnvCI ||
			noticeEnv == env.ServiceEnvTest {
			return nil
		}
		return notice
	})

	return &Notifier{
		Client: client,
	}, nil
}

// Notifier implements diagnostics.Notifier.
type Notifier struct {
	Client *gobrake.Notifier
}

// Notify sends an error.
func (n *Notifier) Notify(err interface{}) error {
	_, sendErr := n.Client.SendNotice(NewNotice(err, nil))
	return ex.New(sendErr)
}

// NotifyWithRequest sends an error with a request.
func (n *Notifier) NotifyWithRequest(err interface{}, req *http.Request) error {
	_, sendErr := n.Client.SendNotice(NewNotice(err, req))
	return ex.New(sendErr)
}

// NotifyWithParams sends an error with custom parameters attached.
func (n *Notifier) NotifyWithParams(err interface{}, params Params) error {
	_, sendErr := n.Client.SendNotice(NewNoticeWithParams(err, params))
	return ex.New(sendErr)
}

func getDefaultContext() map[string]interface{} {
	defaultContextOnce.Do(func() {
		defaultContext = map[string]interface{}{
			"notifier": map[string]interface{}{
				"name":    "gobrake",
				"version": "3.4.0",
				"url":     "https://github.com/airbrake/gobrake",
			},
			"language":     runtime.Version(),
			"os":           runtime.GOOS,
			"architecture": runtime.GOARCH,
		}

		if s, err := os.Hostname(); err == nil {
			defaultContext["hostname"] = s
		}

		if wd, err := os.Getwd(); err == nil {
			defaultContext["rootDirectory"] = wd
		}

		if s := os.Getenv("GOPATH"); s != "" {
			defaultContext["gopath"] = s
		}
	})
	return defaultContext
}
