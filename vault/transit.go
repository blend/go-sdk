package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
)

// Assert Transit implements TransitClient
var (
	_ TransitClient = (*Transit)(nil)
)

// Transit defines vault transit interactions
type Transit struct {
	Client *APIClient
}

// CreateTransitKey creates a transit key path
func (vt Transit) CreateTransitKey(ctx context.Context, key string, options ...CreateTransitKeyOption) error {
	var config CreateTransitKeyConfig
	for _, o := range options {
		err := o(&config)
		if err != nil {
			return err
		}
	}

	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/keys/", key)).WithContext(ctx)

	body, err := vt.Client.jsonBody(config)
	if err != nil {
		return err
	}
	req.Body = body

	res, err := vt.Client.send(req, OptTraceVaultOperation("transit.create"), OptTraceKeyName(key))
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

// ConfigureTransitKey configures a transit key path
func (vt Transit) ConfigureTransitKey(ctx context.Context, key string, options ...UpdateTransitKeyOption) error {
	var config UpdateTransitKeyConfig
	for _, o := range options {
		err := o(&config)
		if err != nil {
			return err
		}
	}

	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/keys/", key, "config")).WithContext(ctx)

	body, err := vt.Client.jsonBody(config)
	if err != nil {
		return err
	}
	req.Body = body

	res, err := vt.Client.send(req, OptTraceVaultOperation("transit.configure"), OptTraceKeyName(key))
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

// ReadTransitKey returns data about a transit key path
func (vt Transit) ReadTransitKey(ctx context.Context, key string) (map[string]interface{}, error) {
	req := vt.Client.createRequest(MethodGet, filepath.Join("/v1/transit/keys/", key)).WithContext(ctx)

	res, err := vt.Client.send(req, OptTraceVaultOperation("transit.read"), OptTraceKeyName(key))
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer res.Close()

	var keyResult TransitKey
	if err = json.NewDecoder(res).Decode(&keyResult); err != nil {
		return nil, err
	}

	return keyResult.Data, nil
}

// DeleteTransitKey deletes a transit key path
func (vt Transit) DeleteTransitKey(ctx context.Context, key string) error {
	req := vt.Client.createRequest(MethodDelete, filepath.Join("/v1/transit/keys/", key)).WithContext(ctx)

	res, err := vt.Client.send(req, OptTraceVaultOperation("transit.delete"), OptTraceKeyName(key))
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

// Encrypt encrypts a given set of data
//
// It is required to create the transit key *before* you use it to encrypt or decrypt data.
func (vt Transit) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
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

	res, err := vt.Client.send(req, OptTraceVaultOperation("transit.encrypt"), OptTraceKeyName(key))
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
//
// It is required to create the transit key *before* you use it to encrypt or decrypt data.
func (vt Transit) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
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

	res, err := vt.Client.send(req, OptTraceVaultOperation("transit.decrypt"), OptTraceKeyName(key))
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
