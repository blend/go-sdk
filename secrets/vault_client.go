package secrets

import (
	"context"
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

// assert VaultClient implements Client
var (
	_ Client = (*VaultClient)(nil)
)

// New creates a new vault client with a default set of options.
func New(options ...Option) (*VaultClient, error) {
	remote, err := url.ParseRequestURI(DefaultAddr)
	if err != nil {
		return nil, err
	}

	client := &VaultClient{
		Timeout:    DefaultTimeout,
		Remote:     remote,
		Mount:      DefaultMount,
		BufferPool: bufferutil.NewPool(DefaultBufferPoolSize),
	}

	client.KV1 = &KV1{Client: client}
	client.KV2 = &KV2{Client: client}
	client.Transit = &VaultTransit{Client: client}

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

// Must does things with the error such as panic.
func Must(c *VaultClient, err error) *VaultClient {
	if err != nil {
		panic(err)
	}
	return c
}

// VaultClient is a client to talk to the secrets store.
type VaultClient struct {
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
	CertPool   *CertPool
	Tracer     Tracer
}

// Put puts a value.
func (c *VaultClient) Put(ctx context.Context, key string, data Values, options ...RequestOption) error {
	backend, err := c.backendKV(ctx, key)
	if err != nil {
		return err
	}

	return backend.Put(ctx, key, data, options...)
}

// Get gets a value at a given key.
func (c *VaultClient) Get(ctx context.Context, key string, options ...RequestOption) (Values, error) {
	backend, err := c.backendKV(ctx, key)
	if err != nil {
		return nil, err
	}

	return backend.Get(ctx, key, options...)
}

// Delete puts a key.
func (c *VaultClient) Delete(ctx context.Context, key string, options ...RequestOption) error {
	backend, err := c.backendKV(ctx, key)
	if err != nil {
		return err
	}
	return backend.Delete(ctx, key, options...)
}

// List returns a slice of key and subfolder names at this path.
func (c *VaultClient) List(ctx context.Context, path string, options ...RequestOption) ([]string, error) {
	backend, err := c.backendKV(ctx, path)
	if err != nil {
		return nil, err
	}

	return backend.List(ctx, path, options...)
}

// ReadInto reads a secret into an object.
func (c *VaultClient) ReadInto(ctx context.Context, key string, obj interface{}, options ...RequestOption) error {
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
func (c *VaultClient) WriteInto(ctx context.Context, key string, obj interface{}, options ...RequestOption) error {
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
func (c *VaultClient) CreateTransitKey(ctx context.Context, key string, options ...CreateTransitKeyOption) error {
	return c.Transit.CreateTransitKey(ctx, key, options...)
}

// ConfigureTransitKey configures a transit key path
func (c *VaultClient) ConfigureTransitKey(ctx context.Context, key string, options ...UpdateTransitKeyOption) error {
	return c.Transit.ConfigureTransitKey(ctx, key, options...)
}

// ReadTransitKey returns data about a transit key path
func (c *VaultClient) ReadTransitKey(ctx context.Context, key string) (map[string]interface{}, error) {
	return c.Transit.ReadTransitKey(ctx, key)
}

// DeleteTransitKey deletes a transit key path
func (c *VaultClient) DeleteTransitKey(ctx context.Context, key string) error {
	return c.Transit.DeleteTransitKey(ctx, key)
}

// Encrypt encrypts a given set of data.
func (c *VaultClient) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
	return c.Transit.Encrypt(ctx, key, context, data)
}

// Decrypt decrypts a given set of data.
func (c *VaultClient) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
	return c.Transit.Decrypt(ctx, key, context, ciphertext)
}

// --------------------------------------------------------------------------------
// utility methods
// --------------------------------------------------------------------------------

func (c *VaultClient) backendKV(ctx context.Context, key string) (KV, error) {
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

func (c *VaultClient) getVersion(ctx context.Context, key string) (string, error) {
	meta, err := c.getMountMeta(ctx, filepath.Join(c.Mount, key))
	if err != nil {
		return "", err
	}
	return meta.Data.Options["version"], nil
}

func (c *VaultClient) getMountMeta(ctx context.Context, key string) (*MountResponse, error) {
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

func (c *VaultClient) jsonBody(input interface{}) (io.ReadCloser, error) {
	buf := c.BufferPool.Get()
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		return nil, err
	}
	return bufferutil.PutOnClose(buf, c.BufferPool), nil
}

func (c *VaultClient) readJSON(r io.Reader, output interface{}) error {
	return json.NewDecoder(r).Decode(output)
}

// copyRemote returns a copy of our remote.
func (c *VaultClient) copyRemote() *url.URL {
	remoteCopy := *c.Remote
	return &remoteCopy
}

// applyOptions applies options to a request.
func (c *VaultClient) applyOptions(req *http.Request, options ...RequestOption) {
	for _, opt := range options {
		opt(req)
	}
}

func (c *VaultClient) createRequest(method, path string, options ...RequestOption) *http.Request {
	remote := c.copyRemote()
	remote.Path = path
	req := &http.Request{
		Method: method,
		URL:    remote,
		Header: http.Header{
			HeaderVaultToken: []string{c.Token},
		},
	}
	c.applyOptions(req, options...)
	return req
}

func (c *VaultClient) send(req *http.Request, traceOptions ...TraceOption) (io.ReadCloser, error) {
	logger.MaybeTrigger(req.Context(), c.Log, NewEvent(req))
	var finisher TraceFinisher
	if c.Tracer != nil {
		var traceErr error
		finisher, traceErr = c.Tracer.Start(req.Context(), traceOptions...)
		if traceErr != nil {
			logger.MaybeError(c.Log, traceErr)
		}
	}
	res, err := c.Client.Do(req)
	if finisher != nil {
		var statusCode = 500
		if res != nil {
			statusCode = res.StatusCode
		}
		finisher.Finish(req.Context(), statusCode, err)
	}
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		buf := c.BufferPool.Get()
		defer c.BufferPool.Put(buf)

		io.Copy(buf, res.Body)
		return nil, ex.New(ExceptionClassForStatus(res.StatusCode), ex.OptMessagef("status: %d; %v", res.StatusCode, buf.String()))
	}
	return res.Body, nil
}

func (c *VaultClient) discard(res io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer res.Close()
	if _, err := io.Copy(ioutil.Discard, res); err != nil {
		return err
	}
	return nil
}
