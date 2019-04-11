package secrets

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
)

// assert VaultTransit implements Transit
var (
	_ Transit = VaultTransit{}
)

// VaultTransit defines vault transit interactions
type VaultTransit struct {
	Client *VaultClient
}

// TransitEncrypt encrypts a given set of data.
func (vt VaultTransit) TransitEncrypt(key string, context map[string]interface{}, data []byte) (string, error) {
	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/encrypt/", key))

	payload := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString(data),
	}
	if context != nil {
		contextJSON, _ := json.Marshal(context)
		contextEncoded := base64.StdEncoding.EncodeToString(contextJSON)
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

// TransitDecrypt decrypts a given set of data.
func (vt VaultTransit) TransitDecrypt(key string, context map[string]interface{}, ciphertext string) ([]byte, error) {
	req := vt.Client.createRequest(MethodPost, filepath.Join("/v1/transit/decrypt/", key))

	payload := map[string]interface{}{
		"ciphertext": ciphertext,
	}
	if context != nil {
		contextJSON, err := json.Marshal(context)
		if err != nil {
			return nil, err
		}
		contextEncoded := base64.StdEncoding.EncodeToString(contextJSON)
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
