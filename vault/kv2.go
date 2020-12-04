package vault

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
)

// assert kv1 implements kv.
var (
	_ KV = (*KV2)(nil)
)

// KV2 defines key value version 2 interactions
type KV2 struct {
	Client *APIClient
}

// Put puts a value.
func (kv2 KV2) Put(ctx context.Context, path string, data Values, options ...CallOption) error {
	contents, err := kv2.Client.jsonBody(SecretData{Data: data})
	if err != nil {
		return err
	}
	req := kv2.Client.createRequest(MethodPut, filepath.Join("/v1/", kv2.fixSecretDataPrefix(path)), options...).WithContext(ctx)
	req.Body = contents
	res, err := kv2.Client.send(req, OptTraceVaultOperation("kv2.put"), OptTraceKeyName(path))
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// Get gets a value at a given key.
func (kv2 KV2) Get(ctx context.Context, path string, options ...CallOption) (Values, error) {
	req := kv2.Client.createRequest(MethodGet, filepath.Join("/v1/", kv2.fixSecretDataPrefix(path)), options...).WithContext(ctx)
	res, err := kv2.Client.send(req, OptTraceVaultOperation("kv2.get"), OptTraceKeyName(path))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response SecretV2
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data.Data, nil
}

// Delete deletes a secret.
func (kv2 KV2) Delete(ctx context.Context, path string, options ...CallOption) error {
	req := kv2.Client.createRequest(MethodDelete, filepath.Join("/v1/", kv2.fixSecretDataPrefix(path)), options...).WithContext(ctx)

	res, err := kv2.Client.send(req, OptTraceVaultOperation("kv2.delete"), OptTraceKeyName(path))
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// List returns a slice of key and subfolder names at this path.
func (kv2 KV2) List(ctx context.Context, path string, options ...CallOption) ([]string, error) {
	req := kv2.Client.createRequest(MethodList, filepath.Join("/v1/", kv2.fixSecretDataPrefix(path)), options...).WithContext(ctx)
	res, err := kv2.Client.send(req, OptTraceVaultOperation("kv2.list"))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response SecretListV2
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data.Keys, nil
}

// fixSecretDataPrefix ensures that a key is prefixed with secret/data/...
func (kv2 KV2) fixSecretDataPrefix(path string) string {
	path = strings.TrimPrefix(path, "/")
	if strings.HasPrefix(path, "secret") && !strings.HasPrefix(path, "secret/data") {
		path = strings.TrimPrefix(path, "secret/")
		path = "secret/data" + path
	}
	return path
}
