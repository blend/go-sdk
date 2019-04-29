package secrets

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
)

// assert VaultTransit implements TransitClient
var (
	_ TransitClient = VaultTransit{}
)

// VaultTransit defines vault transit interactions
type VaultTransit struct {
	Client *VaultClient
}

// CreateTransitKey creates a transit key path
func (vt VaultTransit) CreateTransitKey(ctx context.Context, key string, params map[string]interface{}) error {
	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/keys/", key)).WithContext(ctx)

	if _, ok := params["type"]; !ok {
		params["type"] = "aes256-gcm96"
	}

	if _, ok := params["derived"]; !ok {
		params["derived"] = true
	}

	body, err := vt.Client.jsonBody(params)
	if err != nil {
		return err
	}
	req.Body = body

	res, err := vt.Client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

// ConfigureTransitKey configures a transit key path
func (vt VaultTransit) ConfigureTransitKey(ctx context.Context, key string, config map[string]interface{}) error {
	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/keys/", key, "config")).WithContext(ctx)

	body, err := vt.Client.jsonBody(config)
	if err != nil {
		return err
	}
	req.Body = body

	res, err := vt.Client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

// ReadTransitKey returns data about a transit key path
func (vt VaultTransit) ReadTransitKey(ctx context.Context, key string) (map[string]interface{}, error) {
	req := vt.Client.createRequest(MethodGet, filepath.Join("/v1/transit/keys/", key)).WithContext(ctx)

	res, err := vt.Client.send(req)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer res.Close()

	var keyMetaResult map[string]interface{}
	if err = json.NewDecoder(res).Decode(&keyMetaResult); err != nil {
		return nil, err
	}

	return keyMetaResult, nil
}

// DeleteTransitKey deletes a transit key path
func (vt VaultTransit) DeleteTransitKey(ctx context.Context, key string) error {
	req := vt.Client.createRequest(MethodDelete, filepath.Join("/v1/transit/keys/", key)).WithContext(ctx)

	res, err := vt.Client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

// Encrypt encrypts a given set of data.
func (vt VaultTransit) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/encrypt/", key)).WithContext(ctx)

	payload := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString(data),
	}
	if context != nil {
		contextEncoded := base64.StdEncoding.EncodeToString(context)
		payload["context"] = contextEncoded
	}
	body, err := vt.Client.jsonBody(payload)
	if err != nil {
		return "", err
	}
	req.Body = body

	res, err := vt.Client.send(req)
	if err != nil {
		return "", err
	}
	defer res.Close()

	var encryptionResult TransitResult
	if err = json.NewDecoder(res).Decode(&encryptionResult); err != nil {
		return "", err
	}

	return encryptionResult.Data.Ciphertext, nil
}

// Decrypt decrypts a given set of data.
func (vt VaultTransit) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/decrypt/", key)).WithContext(ctx)

	payload := map[string]interface{}{
		"ciphertext": ciphertext,
	}
	if context != nil {
		contextEncoded := base64.StdEncoding.EncodeToString(context)
		payload["context"] = contextEncoded
	}
	body, err := vt.Client.jsonBody(payload)
	if err != nil {
		return nil, err
	}
	req.Body = body

	res, err := vt.Client.send(req)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var decryptionResult TransitResult
	if err = json.NewDecoder(res).Decode(&decryptionResult); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(decryptionResult.Data.Plaintext)
}
