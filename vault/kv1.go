/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import (
	"context"
	"encoding/json"
	"path/filepath"
)

// assert KV1 implements kv.
var (
	_ KV = &KV1{}
)

// KV1 defines key value version 1 interactions
type KV1 struct {
	Client *APIClient
}

// Put puts a value.
func (kv1 KV1) Put(ctx context.Context, path string, data Values, options ...CallOption) error {
	contents, err := kv1.Client.jsonBody(data)
	if err != nil {
		return err
	}
	req := kv1.Client.createRequest(MethodPut, filepath.Join("/v1/", path), options...).WithContext(ctx)
	req.Body = contents
	res, err := kv1.Client.send(req, OptTraceVaultOperation("kv1.put"), OptTraceKeyName(path))
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// Get gets a value at a given key.
func (kv1 KV1) Get(ctx context.Context, path string, options ...CallOption) (Values, error) {
	req := kv1.Client.createRequest(MethodGet, filepath.Join("/v1/", path), options...).WithContext(ctx)
	res, err := kv1.Client.send(req, OptTraceVaultOperation("kv1.get"), OptTraceKeyName(path))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response SecretV1
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

// Delete puts a key.
func (kv1 KV1) Delete(ctx context.Context, path string, options ...CallOption) error {
	req := kv1.Client.createRequest(MethodDelete, filepath.Join("/v1/", path), options...).WithContext(ctx)
	res, err := kv1.Client.send(req, OptTraceVaultOperation("kv1.delete"), OptTraceKeyName(path))
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// List returns a slice of key and subfolder names at this path.
func (kv1 KV1) List(ctx context.Context, path string, options ...CallOption) ([]string, error) {
	req := kv1.Client.createRequest(MethodList, filepath.Join("/v1/", path), options...).WithContext(ctx)
	res, err := kv1.Client.send(req, OptTraceVaultOperation("kv1.list"))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response SecretListV1
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data.Keys, nil
}
