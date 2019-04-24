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
