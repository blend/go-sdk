package gobrake

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const waitTimeout = 5 * time.Second

const httpEnhanceYourCalm = 420
const httpStatusTooManyRequests = 429

const maxNoticeLen = 64 * 1024

var (
	errClosed             = errors.New("gobrake: notifier is closed")
	errQueueFull          = errors.New("gobrake: queue is full (error is dropped)")
	errUnauthorized       = errors.New("gobrake: unauthorized: invalid project id or key")
	errAccountRateLimited = errors.New("gobrake: account is rate limited")
	errIPRateLimited      = errors.New("gobrake: IP is rate limited")
	errNoticeTooBig       = errors.New("gobrake: notice exceeds 64KB max size limit")
)

var (
	httpClientOnce sync.Once
	httpClient     *http.Client
)

func defaultHTTPClient() *http.Client {
	httpClientOnce.Do(func() {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout:   15 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
				TLSClientConfig: &tls.Config{
					ClientSessionCache: tls.NewLRUClientSessionCache(1024),
				},
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: 10 * time.Second,
			},
			Timeout: 10 * time.Second,
		}
	})
	return httpClient
}

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type filter func(*Notice) *Notice

type NotifierOptions struct {
	// Airbrake project id.
	ProjectId int64
	// Airbrake project key.
	ProjectKey string
	// Airbrake host name. Default is https://airbrake.io.
	Host string

	// Environment such as production or development.
	Environment string
	// Git revision. Default is SOURCE_VERSION on Heroku.
	Revision string
	// List of keys containing sensitive information that must be filtered out.
	// Default is password, secret.
	KeysBlacklist []interface{}

	// http.Client that is used to interact with Airbrake API.
	HTTPClient *http.Client
}

func (opt *NotifierOptions) init() {
	if opt.Host == "" {
		opt.Host = "https://api.airbrake.io"
	}

	if opt.Revision == "" {
		// https://devcenter.heroku.com/changelog-items/630
		opt.Revision = os.Getenv("SOURCE_VERSION")
	}

	if opt.KeysBlacklist == nil {
		opt.KeysBlacklist = []interface{}{
			regexp.MustCompile("password"),
			regexp.MustCompile("secret"),
		}
	}

	if opt.HTTPClient == nil {
		opt.HTTPClient = defaultHTTPClient()
	}
}

type Notifier struct {
	opt             *NotifierOptions
	createNoticeURL string

	filters []filter

	inFlight int32 // atomic
	limit    chan struct{}
	wg       sync.WaitGroup

	Routes  *routeStats
	Queries *QueryStats

	rateLimitReset uint32 // atomic
	_closed        uint32 // atomic
}

func NewNotifierWithOptions(opt *NotifierOptions) *Notifier {
	opt.init()

	n := &Notifier{
		opt: opt,
		createNoticeURL: fmt.Sprintf("%s/api/v3/projects/%d/notices",
			opt.Host, opt.ProjectId),

		limit: make(chan struct{}, 2*runtime.NumCPU()),

		Routes:  newRouteStats(opt),
		Queries: newQueryStats(opt),
	}

	n.AddFilter(newNotifierFilter(n))
	n.AddFilter(gopathFilter)
	n.AddFilter(gitFilter)

	if len(opt.KeysBlacklist) > 0 {
		n.AddFilter(NewBlacklistKeysFilter(opt.KeysBlacklist...))
	}

	return n
}

func NewNotifier(projectId int64, projectKey string) *Notifier {
	return NewNotifierWithOptions(&NotifierOptions{
		ProjectId:  projectId,
		ProjectKey: projectKey,
	})
}

// AddFilter adds filter that can change notice or ignore it by returning nil.
func (n *Notifier) AddFilter(fn func(*Notice) *Notice) {
	n.filters = append(n.filters, fn)
}

// Notify notifies Airbrake about the error.
func (n *Notifier) Notify(e interface{}, req *http.Request) {
	notice := n.Notice(e, req, 1)
	n.SendNoticeAsync(notice)
}

// Notice returns Aibrake notice created from error and request. depth
// determines which call frame to use when constructing backtrace.
func (n *Notifier) Notice(err interface{}, req *http.Request, depth int) *Notice {
	return NewNotice(err, req, depth+3)
}

type sendResponse struct {
	Id string `json:"id"`
}

// SendNotice sends notice to Airbrake.
func (n *Notifier) SendNotice(notice *Notice) (string, error) {
	if n.closed() {
		return "", errClosed
	}
	return n.sendNotice(notice)
}

func (n *Notifier) sendNotice(notice *Notice) (string, error) {
	for _, fn := range n.filters {
		notice = fn(notice)
		if notice == nil {
			// Notice is ignored.
			return "", nil
		}
	}

	if time.Now().Unix() < int64(atomic.LoadUint32(&n.rateLimitReset)) {
		return "", errIPRateLimited
	}

	buf := buffers.Get().(*bytes.Buffer)
	defer buffers.Put(buf)

	buf.Reset()
	err := json.NewEncoder(buf).Encode(notice)
	if err != nil {
		return "", err
	}

	if buf.Len() > maxNoticeLen {
		return "", errNoticeTooBig
	}

	req, err := http.NewRequest("POST", n.createNoticeURL, buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+n.opt.ProjectKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.opt.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var sendResp sendResponse
		err = json.NewDecoder(buf).Decode(&sendResp)
		if err != nil {
			return "", err
		}
		return sendResp.Id, nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return "", errUnauthorized
	case httpStatusTooManyRequests:
		delayStr := resp.Header.Get("X-RateLimit-Delay")
		delay, err := strconv.ParseInt(delayStr, 10, 64)
		if err == nil {
			atomic.StoreUint32(&n.rateLimitReset, uint32(time.Now().Unix()+delay))
		}
		return "", errIPRateLimited
	case httpEnhanceYourCalm:
		return "", errAccountRateLimited
	}

	err = fmt.Errorf("got unexpected response status=%q", resp.Status)
	logger.Printf("SendNotice failed reporting notice=%q: %s", notice, err)
	return "", err
}

// SendNoticeAsync is like SendNotice, but sends notice asynchronously.
// Pending notices can be flushed with Flush.
func (n *Notifier) SendNoticeAsync(notice *Notice) {
	if n.closed() {
		notice.Error = errClosed
		return
	}

	inFlight := atomic.AddInt32(&n.inFlight, 1)
	if inFlight > 1000 {
		atomic.AddInt32(&n.inFlight, -1)
		notice.Error = errQueueFull
		return
	}

	n.wg.Add(1)
	go func() {
		n.limit <- struct{}{}

		notice.Id, notice.Error = n.sendNotice(notice)
		atomic.AddInt32(&n.inFlight, -1)
		n.wg.Done()

		<-n.limit
	}()
}

// NotifyOnPanic notifies Airbrake about the panic and should be used
// with defer statement.
func (n *Notifier) NotifyOnPanic() {
	if v := recover(); v != nil {
		notice := n.Notice(v, nil, 3)
		notice.Context["severity"] = "critical"
		n.SendNotice(notice)
		panic(v)
	}
}

// Flush waits for pending requests to finish.
func (n *Notifier) Flush() {
	n.waitTimeout(waitTimeout)
}

func (n *Notifier) Close() error {
	return n.CloseTimeout(waitTimeout)
}

// CloseTimeout waits for pending requests to finish and then closes the notifier.
func (n *Notifier) CloseTimeout(timeout time.Duration) error {
	if !atomic.CompareAndSwapUint32(&n._closed, 0, 1) {
		return nil
	}
	return n.waitTimeout(timeout)
}

func (n *Notifier) closed() bool {
	return atomic.LoadUint32(&n._closed) == 1
}

func (n *Notifier) waitTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Wait timed out after %s", timeout)
	}
}
