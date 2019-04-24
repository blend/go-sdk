package secrets

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
)

// assert VaultClient implements Client
var _ Client = &VaultClient{}

// assert VaultClient implements TransitClient
var _ TransitClient = &VaultClient{}

// NewVaultClient returns a new client.
func NewVaultClient() (*VaultClient, error) {
	return NewVaultClientFromConfig(&Config{})
}

// NewVaultClientFromConfig returns a new client from a config.
func NewVaultClientFromConfig(cfg *Config) (*VaultClient, error) {
	xport := &http.Transport{}
	err := http2.ConfigureTransport(xport)
	if err != nil {
		return nil, err
	}
	remote, err := url.ParseRequestURI(cfg.GetAddr())
	if err != nil {
		return nil, err
	}
	var certPool *CertPool
	if caPaths := cfg.GetRootCAs(); len(caPaths) > 0 {
		certPool, err = NewCertPool()
		if err != nil {
			return nil, err
		}
		err = certPool.AddPaths(caPaths...)
		if err != nil {
			return nil, err
		}
		xport.TLSClientConfig = &tls.Config{
			RootCAs: certPool.Pool(),
		}
	}
	client := &VaultClient{
		remote:     remote,
		mount:      cfg.GetMount(),
		bufferPool: NewBufferPool(DefaultBufferPoolSize),
		token:      cfg.GetToken(),
		certPool:   certPool,
		client: &http.Client{
			Timeout:   cfg.GetTimeout(),
			Transport: xport,
		},
	}

	client.kv1 = &kv1{client: client}
	client.kv2 = &kv2{client: client}
	client.Transit = &VaultTransit{Client: client}

	return client, nil
}

// NewVaultClientFromEnv is a helper to create a client from a config read from the environment.
func NewVaultClientFromEnv() (*VaultClient, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewVaultClientFromConfig(cfg)
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
	remote *url.URL
	token  string
	mount  string
	log    *logger.Logger

	kv1     *kv1
	kv2     *kv2
	Transit TransitClient

	bufferPool *BufferPool
	client     HTTPClient
	certPool   *CertPool
}

// WithRemote set the client remote url.
func (c *VaultClient) WithRemote(remote *url.URL) *VaultClient {
	c.remote = remote
	return c
}

// Remote returns the client remote addr.
func (c *VaultClient) Remote() *url.URL {
	return c.remote
}

// WithToken sets the token.
func (c *VaultClient) WithToken(token string) *VaultClient {
	c.token = token
	return c
}

// Token returns the token.
func (c *VaultClient) Token() string {
	return c.token
}

// WithMount sets the token.
func (c *VaultClient) WithMount(mount string) *VaultClient {
	c.mount = mount
	return c
}

// Mount returns the mount.
func (c *VaultClient) Mount() string {
	return c.mount
}

// WithHTTPClient sets the http client.
func (c *VaultClient) WithHTTPClient(hc HTTPClient) *VaultClient {
	c.client = hc
	return c
}

// HTTPClient sets the http client.
func (c *VaultClient) HTTPClient() HTTPClient {
	return c.client
}

// CertPool returns the cert pool.
func (c *VaultClient) CertPool() *CertPool {
	return c.certPool
}

// WithLogger sets the logger.
func (c *VaultClient) WithLogger(log *logger.Logger) *VaultClient {
	c.log = log
	return c
}

// Logger returns the logger.
func (c *VaultClient) Logger() *logger.Logger {
	return c.log
}

// Put puts a value.
func (c *VaultClient) Put(key string, data Values, options ...Option) error {
	backend, err := c.backendKV(key)
	if err != nil {
		return err
	}

	return backend.Put(key, data, options...)
}

// Get gets a value at a given key.
func (c *VaultClient) Get(key string, options ...Option) (Values, error) {
	backend, err := c.backendKV(key)
	if err != nil {
		return nil, err
	}

	return backend.Get(key, options...)
}

// Delete puts a key.
func (c *VaultClient) Delete(key string, options ...Option) error {
	backend, err := c.backendKV(key)
	if err != nil {
		return err
	}
	return backend.Delete(key, options...)
}

// ReadInto reads a secret into an object.
func (c *VaultClient) ReadInto(key string, obj interface{}, options ...Option) error {
	response, err := c.Get(key, options...)
	if err != nil {
		return err
	}
	return util.Reflection.PatchStrings(ReflectTagName, response, obj)
}

// WriteInto writes an object into a secret at a given key.
func (c *VaultClient) WriteInto(key string, obj interface{}, options ...Option) error {
	return c.Put(key, util.Reflection.DecomposeStrings(obj, ReflectTagName), options...)
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

func (c *VaultClient) backendKV(key string) (KV, error) {
	version, err := c.getVersion(key)
	if err != nil {
		return nil, err
	}
	switch version {
	case Version1:
		return c.kv1, nil
	case Version2:
		return c.kv2, nil
	default:
		return c.kv1, nil
	}
}

func (c *VaultClient) getVersion(key string) (string, error) {
	meta, err := c.getMountMeta(filepath.Join(c.mount, key))
	if err != nil {
		return "", err
	}
	return meta.Data.Options["version"], nil
}

func (c *VaultClient) getMountMeta(key string) (*MountResponse, error) {
	req := c.createRequest(MethodGet, filepath.Join("/v1/sys/internal/ui/mounts/", key))

	res, err := c.client.Do(req)
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
	buf := c.bufferPool.Get()
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *VaultClient) readJSON(r io.Reader, output interface{}) error {
	return json.NewDecoder(r).Decode(output)
}

// copyRemote returns a copy of our remote.
func (c *VaultClient) copyRemote() *url.URL {
	remoteCopy := *c.remote
	return &remoteCopy
}

// applyOptions applies options to a request.
func (c *VaultClient) applyOptions(req *http.Request, options ...Option) {
	for _, opt := range options {
		opt(req)
	}
}

func (c *VaultClient) createRequest(method, path string, options ...Option) *http.Request {
	remote := c.copyRemote()
	remote.Path = path
	req := &http.Request{
		Method: method,
		URL:    remote,
		Header: http.Header{
			HeaderVaultToken: []string{c.Token()},
		},
	}
	c.applyOptions(req, options...)
	return req
}

func (c *VaultClient) send(req *http.Request) (io.ReadCloser, error) {
	if c.log != nil {
		c.log.Trigger(NewEvent(req))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		buf := c.bufferPool.Get()
		defer buf.Close()
		io.Copy(buf, res.Body)
		return nil, fmt.Errorf("non-2xx returned from remote: %d; %v", res.StatusCode, buf.String())
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
