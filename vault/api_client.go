package vault

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/bufferutil"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// Assert APIClient implements Client
var (
	_ Client = (*APIClient)(nil)
)

// New creates a new vault client with a default set of options.
func New(options ...Option) (*APIClient, error) {
	remote, err := url.ParseRequestURI(DefaultAddr)
	if err != nil {
		return nil, err
	}

	client := &APIClient{
		Timeout:    DefaultTimeout,
		Mount:      DefaultMount,
		Remote:     remote,
		BufferPool: bufferutil.NewPool(DefaultBufferPoolSize),
	}

	client.KV1 = &KV1{Client: client}
	client.KV2 = &KV2{Client: client}
	client.Transit = &Transit{Client: client}

	for _, option := range options {
		if err = option(client); err != nil {
			return nil, err
		}
	}

	xport := client.Transport
	if xport == nil {
		xport = &http.Transport{}
		err = http2.ConfigureTransport(xport)
		if err != nil {
			return nil, err
		}
	}

	client.Client = &http.Client{
		Transport: xport,
		Timeout:   client.Timeout,
	}

	return client, nil
}

// APIClient is a client to talk to vault.
type APIClient struct {
	Timeout    time.Duration
	Transport  *http.Transport
	Remote     *url.URL
	Token      string
	Mount      string
	Log        logger.Log
	BufferPool *bufferutil.Pool
	KV1        *KV1
	KV2        *KV2
	Transit    TransitClient
	Client     HTTPClient
	CertPool   *x509.CertPool
	Tracer     Tracer
}

// Put puts a value.
func (c *APIClient) Put(ctx context.Context, key string, data Values, options ...CallOption) error {
	backend, err := c.backendKV(ctx, key)
	if err != nil {
		return err
	}

	return backend.Put(ctx, key, data, options...)
}

// Get gets a value at a given key.
func (c *APIClient) Get(ctx context.Context, key string, options ...CallOption) (Values, error) {
	backend, err := c.backendKV(ctx, key)
	if err != nil {
		return nil, err
	}

	return backend.Get(ctx, key, options...)
}

// Delete puts a key.
func (c *APIClient) Delete(ctx context.Context, key string, options ...CallOption) error {
	backend, err := c.backendKV(ctx, key)
	if err != nil {
		return err
	}
	return backend.Delete(ctx, key, options...)
}

// List returns a slice of key and subfolder names at this path.
func (c *APIClient) List(ctx context.Context, path string, options ...CallOption) ([]string, error) {
	backend, err := c.backendKV(ctx, path)
	if err != nil {
		return nil, err
	}

	return backend.List(ctx, path, options...)
}

// ReadInto reads a secret into an object.
func (c *APIClient) ReadInto(ctx context.Context, key string, obj interface{}, options ...CallOption) error {
	response, err := c.Get(ctx, key, options...)
	if err != nil {
		return err
	}
	asStrings := make(map[string]string)
	for k, v := range response {
		if s, ok := v.(string); ok {
			asStrings[k] = s
		}
	}
	return RestoreJSON(asStrings, obj)
}

// WriteInto writes an object into a secret at a given key.
func (c *APIClient) WriteInto(ctx context.Context, key string, obj interface{}, options ...CallOption) error {
	data, err := DecomposeJSON(obj)
	if err != nil {
		return err
	}
	asData := make(map[string]interface{})
	for k, v := range data {
		asData[k] = v
	}
	return c.Put(ctx, key, asData, options...)
}

// CreateTransitKey creates a transit key path
func (c *APIClient) CreateTransitKey(ctx context.Context, key string, options ...CreateTransitKeyOption) error {
	return c.Transit.CreateTransitKey(ctx, key, options...)
}

// ConfigureTransitKey configures a transit key path
func (c *APIClient) ConfigureTransitKey(ctx context.Context, key string, options ...UpdateTransitKeyOption) error {
	return c.Transit.ConfigureTransitKey(ctx, key, options...)
}

// ReadTransitKey returns data about a transit key path
func (c *APIClient) ReadTransitKey(ctx context.Context, key string) (map[string]interface{}, error) {
	return c.Transit.ReadTransitKey(ctx, key)
}

// DeleteTransitKey deletes a transit key path
func (c *APIClient) DeleteTransitKey(ctx context.Context, key string) error {
	return c.Transit.DeleteTransitKey(ctx, key)
}

// Encrypt encrypts a given set of data.
func (c *APIClient) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
	return c.Transit.Encrypt(ctx, key, context, data)
}

// Decrypt decrypts a given set of data.
func (c *APIClient) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
	return c.Transit.Decrypt(ctx, key, context, ciphertext)
}

// --------------------------------------------------------------------------------
// utility methods
// --------------------------------------------------------------------------------

func (c *APIClient) backendKV(ctx context.Context, key string) (KV, error) {
	version, err := c.getVersion(ctx, key)
	if err != nil {
		return nil, err
	}
	switch version {
	case Version1:
		return c.KV1, nil
	case Version2:
		return c.KV2, nil
	default:
		return c.KV1, nil
	}
}

func (c *APIClient) getVersion(ctx context.Context, key string) (string, error) {
	meta, err := c.getMountMeta(ctx, filepath.Join(c.Mount, key))
	if err != nil {
		return "", err
	}
	return meta.Data.Options["version"], nil
}

func (c *APIClient) getMountMeta(ctx context.Context, key string) (*MountResponse, error) {
	req := c.createRequest(MethodGet, filepath.Join("/v1/sys/internal/ui/mounts/", key))
	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response MountResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *APIClient) jsonBody(input interface{}) (io.ReadCloser, error) {
	buf := c.BufferPool.Get()
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		return nil, err
	}
	return bufferutil.PutOnClose(buf, c.BufferPool), nil
}

func (c *APIClient) readJSON(r io.Reader, output interface{}) error {
	return json.NewDecoder(r).Decode(output)
}

// copyRemote returns a copy of our remote.
func (c *APIClient) copyRemote() *url.URL {
	remoteCopy := *c.Remote
	return &remoteCopy
}

// applyOptions applies options to a request.
func (c *APIClient) applyOptions(req *http.Request, options ...CallOption) error {
	var err error
	for _, opt := range options {
		if err = opt(req); err != nil {
			return err
		}
	}
	return nil
}

func (c *APIClient) createRequest(method, path string, options ...CallOption) *http.Request {
	remote := c.copyRemote()
	remote.Path = path
	req := &http.Request{
		Method: method,
		URL:    remote,
		Header: http.Header{
			HeaderVaultToken: []string{c.Token},
		},
	}
	_ = c.applyOptions(req, options...)
	return req
}

func (c *APIClient) send(req *http.Request, traceOptions ...TraceOption) (body io.ReadCloser, err error) {
	var statusCode int
	var finisher TraceFinisher
	if c.Log != nil {
		e := NewEvent(req)
		start := time.Now()
		defer func() {
			e.Elapsed = time.Since(start)
			logger.MaybeTriggerContext(req.Context(), c.Log, e)
		}()
	}
	if finisher != nil {
		defer func() {
			finisher.Finish(req.Context(), statusCode, err)
		}()
	}
	if c.Tracer != nil {
		var traceErr error
		finisher, traceErr = c.Tracer.Start(req.Context(), traceOptions...)
		if traceErr != nil {
			logger.MaybeError(c.Log, traceErr)
		}
	}

	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		statusCode = 500
		return
	}
	statusCode = res.StatusCode

	if statusCode > 299 {
		buf := c.BufferPool.Get()
		defer c.BufferPool.Put(buf)
		if _, err = io.Copy(buf, res.Body); err != nil {
			err = ex.New(err)
			return
		}
		err = ex.New(ErrClassForStatus(statusCode), ex.OptMessagef("status: %d; %v", statusCode, buf.String()))
		return
	}
	body = res.Body
	return
}

func (c *APIClient) discard(res io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer res.Close()
	if _, err = io.Copy(ioutil.Discard, res); err != nil {
		return err
	}
	return nil
}
