package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/crypto"
)

var (
	_ Client = (*MockClient)(nil)
)

// MockTransitKey is a mocked transit key
type MockTransitKey struct {
	Properties map[string]interface{}
	Keys       map[string][]byte
}

// NewMockClient creates a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{
		SecretValues: make(map[string]Values),
		TransitKeys:  make(map[string]MockTransitKey),
	}
}

// MockClient is a mock events client
type MockClient struct {
	SecretValues map[string]Values
	TransitKeys  map[string]MockTransitKey
}

// Put puts a value.
func (c *MockClient) Put(_ context.Context, key string, data Values, options ...RequestOption) error {
	c.SecretValues[key] = data
	return nil
}

// Get gets a value at a given key.
func (c *MockClient) Get(_ context.Context, key string, options ...RequestOption) (Values, error) {
	val, exists := c.SecretValues[key]
	if !exists {
		return nil, fmt.Errorf("Key not found: %s", key)
	}

	return val, nil
}

// Delete deletes a key.
func (c *MockClient) Delete(_ context.Context, key string, options ...RequestOption) error {
	if _, exists := c.SecretValues[key]; !exists {
		return fmt.Errorf("Key not found: %s", key)
	}

	delete(c.SecretValues, key)
	return nil
}

// List lists keys on a path
func (c *MockClient) List(_ context.Context, path string, options ...RequestOption) ([]string, error) {
	keys := make([]string, 0)
	folderSet := make(map[string]struct{})
	p := path
	if !strings.HasSuffix(path, "/") {
		p = path + "/"
	}
	for k := range c.SecretValues {
		if strings.HasPrefix(k, p) {
			s := strings.TrimPrefix(k, p)
			if strings.ContainsRune(s, '/') {
				folder := fmt.Sprintf("%s/", strings.Split(s, "/")[0])
				if _, ok := folderSet[folder]; !ok {
					folderSet[folder] = struct{}{}
					keys = append(keys, folder)
				}
			} else {
				keys = append(keys, s)
			}
		}
	}
	return keys, nil
}

// CreateTransitKey creates a new transit key.
func (c *MockClient) CreateTransitKey(ctx context.Context, key string, params map[string]interface{}) error {
	c.TransitKeys[key] = MockTransitKey{
		Properties: make(map[string]interface{}),
		Keys:       make(map[string][]byte),
	}

	return nil
}

// ConfigureTransitKey configures a transit key path
func (c *MockClient) ConfigureTransitKey(ctx context.Context, key string, config map[string]interface{}) error {
	keyPath, ok := c.TransitKeys[key]
	if !ok {
		return fmt.Errorf("No key")
	}

	for opt, value := range config {
		keyPath.Properties[opt] = value
	}

	c.TransitKeys[key] = keyPath
	return nil
}

// ReadTransitKey returns data about a transit key path
func (c *MockClient) ReadTransitKey(ctx context.Context, key string) (map[string]interface{}, error) {
	keyPath, ok := c.TransitKeys[key]
	if !ok {
		return map[string]interface{}{}, fmt.Errorf("No key")
	}

	return keyPath.Properties, nil
}

// DeleteTransitKey deletes a transit key path
func (c *MockClient) DeleteTransitKey(ctx context.Context, key string) error {
	keyPath, ok := c.TransitKeys[key]
	if !ok {
		return fmt.Errorf("No key")
	}

	if keyPath.Properties["deletion_allowed"] != true {
		return fmt.Errorf("Deletion is not allowed for key")
	}

	delete(c.TransitKeys, key)
	return nil
}

func (c *MockClient) deriveTransitKey(name string, context []byte) ([]byte, error) {
	contextStr := string(context)

	keyPath, ok := c.TransitKeys[name]
	if !ok {
		return nil, fmt.Errorf("No key")
	}

	key, ok := keyPath.Keys[contextStr]
	if !ok {
		key, _ = crypto.CreateKey(32)
		c.TransitKeys[name].Keys[contextStr] = key
	}

	return key, nil
}

// Encrypt encrypts a given set of data.
func (c *MockClient) Encrypt(ctx context.Context, name string, context, data []byte) (string, error) {
	key, err := c.deriveTransitKey(name, context)
	if err != nil {
		return "", err
	}

	encryptedData, err := crypto.Encrypt(key, data)
	return string(encryptedData), err
}

// Decrypt decrypts a given set of data.
func (c *MockClient) Decrypt(ctx context.Context, name string, context []byte, ciphertext string) ([]byte, error) {
	key, err := c.deriveTransitKey(name, context)
	if err != nil {
		return nil, err
	}

	return crypto.Decrypt(key, []byte(ciphertext))
}
